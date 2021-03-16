package revip

import (
	"context"
	"reflect"
	"strings"
	"time"

	etcd "go.etcd.io/etcd/clientv3"
)

func FromEtcd(client *etcd.Client, namespace string, f Unmarshaler) Option {
	prefix := []string{namespace}

	return func(c Config, m ...OptionMeta) error {
		var ctx context.Context

		for _, mm := range m {
			switch v := mm.(type) {
			case context.Context:
				ctx = v
			}
		}

		if ctx == nil {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(
				context.Background(),
				60*time.Second,
			)
			defer cancel()
		}

		return walkStruct(c, func(v reflect.Value, path []string) error {
			if !v.CanAddr() {
				return skipBranch
			}
			if v.Type().Kind() == reflect.Ptr {
				return nil
			}

			key := strings.Join(append(prefix, path...), PathDelimiter)

			r, err := client.Get(ctx, key)
			if err != nil {
				return err
			}

			for _, kv := range r.Kvs {
				err = f(kv.Value, v.Addr().Interface())
				if err != nil {
					return &ErrUnmarshal{At: key, Err: err}
				}
			}

			return nil
		})
	}
}
