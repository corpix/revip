package revip

import (
	"context"
	"reflect"
	"strings"
	"time"

	etcd "go.etcd.io/etcd/clientv3"
)

// optional context could be providen through meta options
// if not providen then default context will be created with 60s timeout
func ToEtcd(client *etcd.Client, namespace string, f Marshaler) Option {
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
			if v.Type().Kind() == reflect.Ptr {
				return skipBranch
			}

			key := strings.Join(append(prefix, path...), PathDelimiter)

			buf, err := f(v.Interface())
			if err != nil {
				return &ErrMarshal{At: key, Err: err}
			}

			_, err = client.Put(ctx, key, string(buf))
			return err
		})
	}
}