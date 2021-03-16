package revip

import (
	"io"
	"os"
	"reflect"

	json "encoding/json"

	yaml "github.com/go-yaml/yaml"
	toml "github.com/pelletier/go-toml"
)

// Marshaler describes a generic marshal interface for data encoding
// which could be used to extend supported formats by defining new `Option`
// implementations.
type Marshaler = func(v interface{}) ([]byte, error)

var (
	JsonMarshaler Marshaler = json.Marshal
	YamlMarshaler Marshaler = yaml.Marshal
	TomlMarshaler Marshaler = toml.Marshal
)

// ToWriter is an `Option` constructor which creates a thunk
// to write configuration to `r` and encode it with `f` marshaler.
func ToWriter(w io.Writer, f Marshaler) Option {
	return func(c Config, m ...OptionMeta) error {
		err := expectKind(reflect.TypeOf(c), reflect.Ptr)
		if err != nil {
			return err
		}

		buf, err := f(c)
		if err != nil {
			return err
		}

		_, err = w.Write(buf)
		return err
	}
}

// ToFile is an `Option` constructor which creates a thunk
// to write configuration to file addressable by `path` with
// content encoded with `f` marshaler.
func ToFile(path string, f Marshaler) Option {
	return func(c Config, m ...OptionMeta) error {
		err := expectKind(reflect.TypeOf(c), reflect.Ptr)
		if err != nil {
			return err
		}

		r, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0700)
		if err != nil {
			return err
		}
		defer r.Close()

		return ToWriter(r, f)(c, m...)
	}
}
