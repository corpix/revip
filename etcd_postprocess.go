package revip

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	etcd "go.etcd.io/etcd/clientv3"
)

//

type UpdateFromEtcdConfig struct {
	Ctx           context.Context
	BatchSize     int
	BatchDuration time.Duration
	OnError       func(error)
}

type UpdateFromEtcdOption = func(*UpdateFromEtcdConfig)

func UpdatesFromEtcdContext(ctx context.Context) UpdateFromEtcdOption {
	return func(c *UpdateFromEtcdConfig) { c.Ctx = ctx }
}

func UpdatesFromEtcdBatch(size int, duration time.Duration) UpdateFromEtcdOption {
	return func(c *UpdateFromEtcdConfig) {
		c.BatchSize = size
		c.BatchDuration = duration
	}
}

func UpdatesFromEtcdErrorHandler(cb func(error)) UpdateFromEtcdOption {
	return func(c *UpdateFromEtcdConfig) {
		c.OnError = cb
	}
}

//

func etcdWatch(ctx context.Context, c Config, namespace string, client *etcd.Client) ([]etcd.WatchChan, error) {
	chs := []etcd.WatchChan{}
	err := walkStruct(c, func(v reflect.Value, path []string) error {
		switch v.Type().Kind() { // ignore nil's and substructs (substructs may have their own handlers)
		case reflect.Struct:
			return skipBranch
		case reflect.Ptr:
			return nil
		}

		ch := client.Watch(
			ctx,
			strings.Join(
				prefixPath(namespace, path),
				PathDelimiter,
			),
		)
		chs = append(chs, ch)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return chs, nil
}

func etcdWatchPump(ctx context.Context, chs []etcd.WatchChan, events chan etcdUpdateEvent) {
	cases := make([]reflect.SelectCase, len(chs)+1)
	cases[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	}
	for n, ch := range chs {
		cases[n+1] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}
	for {
		chosen, value, ok := reflect.Select(cases)
		if chosen == 0 {
			return // ctx done
		}

		if ok {
			r := value.Interface().(etcd.WatchResponse)
			for _, evt := range r.Events {
				events <- etcdUpdateEvent{
					operation: int32(evt.Type),
					data:      evt.Kv.Value,
					key:       string(evt.Kv.Key),
					version:   evt.Kv.Version,
				}
			}
		}
	}
}

func etcdWatchHandle(ctx context.Context, batchSize int, batchDuration time.Duration, c Config, namespace string, events chan etcdUpdateEvent, v Updateable, onError func(error), f Unmarshaler) {
	var (
		versions = map[string]int64{}
		acc      = make([]etcdUpdateEvent, batchSize)
		n        = 0
		err      error
	)

	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-events:
			if !ok {
				return
			}
			if ver, ok := versions[evt.key]; ok && ver >= evt.version {
				continue // skip already updated keys
			}

			acc[n] = evt
			n++

			if n >= batchSize {
				goto flush
			}
		case <-time.After(batchDuration):
			goto flush
		}
		continue

	flush:
		if n == 0 {
			continue
		}

		dst := reflect.New(reflect.ValueOf(c).Elem().Type()).Interface()
		err = mapstructure.Decode(c, &dst)
		if err != nil {
			onError(err)
			continue
		}

		evtByKey := map[string]etcdUpdateEvent{}
		for _, evt := range acc[:n] {
			evtByKey[evt.key] = evt
		}

		err = walkStruct(dst, func(v reflect.Value, path []string) error {
			k := v.Type().Kind()
			switch { // ignore nil's and substructs (substructs may have their own handlers)
			case k == reflect.Struct:
				return skipBranch
			case k == reflect.Ptr:
				return nil
			case !v.CanAddr():
				return skipBranch
			}

			key := strings.Join(prefixPath(namespace, path), PathDelimiter)
			if evt, ok := evtByKey[key]; ok {
				switch evt.operation {
				case etcdOperationDelete:
					v.Set(reflect.New(v.Type()).Elem())
				case etcdOperationPut:
					err := f(evt.data, v.Addr().Interface())
					if err != nil {
						return err
					}
				}
			}

			return nil
		})
		if err != nil {
			onError(err)
			continue
		}

		// update configuration

		v.Update(dst)

		err = mapstructure.Decode(dst, c)
		if err != nil {
			onError(err)
			continue
		}

		//

		n = 0
	}
}

func WithUpdatesFromEtcd(client *etcd.Client, namespace string, f Unmarshaler, op ...UpdateFromEtcdOption) Option {
	cfg := &UpdateFromEtcdConfig{
		Ctx:           context.Background(),
		BatchSize:     64,
		BatchDuration: 1 * time.Second,
		OnError:       func(err error) { panic(err) },
	}
	for _, apply := range op {
		apply(cfg)
	}

	return func(c Config, m ...OptionMeta) error {
		v, ok := c.(Updateable)
		if !ok {
			return nil
		}

		events := make(chan etcdUpdateEvent, 16)

		chs, err := etcdWatch(cfg.Ctx, c, namespace, client)
		if err != nil {
			return err
		}

		go etcdWatchPump(cfg.Ctx, chs, events)
		go etcdWatchHandle(cfg.Ctx, cfg.BatchSize, cfg.BatchDuration, c, namespace, events, v, cfg.OnError, f)

		return nil
	}
}
