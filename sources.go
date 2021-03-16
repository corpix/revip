package revip

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"syscall"

	json "encoding/json"
	yaml "github.com/go-yaml/yaml"
	env "github.com/kelseyhightower/envconfig"
	toml "github.com/pelletier/go-toml"
)

// Unmarshaler describes a generic unmarshal interface for data decoding
// which could be used to extend supported formats by defining new `Option`
// implementations.
type Unmarshaler = func(in []byte, v interface{}) error

var (
	JsonUnmarshaler Unmarshaler = json.Unmarshal
	YamlUnmarshaler Unmarshaler = yaml.Unmarshal
	TomlUnmarshaler Unmarshaler = toml.Unmarshal
)

// FromReader is an `Option` constructor which creates a thunk
// to read configuration from `r` and decode it with `f` unmarshaler.
// Current implementation buffers all data in memory.
func FromReader(r io.Reader, f Unmarshaler) Option {
	return func(c Config, m ...OptionMeta) error {
		err := expectKind(reflect.TypeOf(c), reflect.Ptr)
		if err != nil {
			return err
		}

		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		return f(buf, c)
	}
}

// FromFile is an `Option` constructor which creates a thunk
// to read configuration from file addressable by `path` with
// conmtent decoded with `f` unmarshaler.
func FromFile(path string, f Unmarshaler) Option {
	return func(c Config, m ...OptionMeta) error {
		err := expectKind(reflect.TypeOf(c), reflect.Ptr)
		if err != nil {
			return err
		}

		r, err := os.Open(path)
		switch e := err.(type) {
		case *os.PathError:
			if e.Err == syscall.ENOENT {
				return &ErrFileNotFound{
					Path: path,
					Err:  err,
				}
			}
		case nil:
		default:
			return err
		}
		defer r.Close()

		return FromReader(r, f)(c, m...)
	}
}

// FromEnviron is an `Option` constructor which creates a thunk
// to read configuration from environment.
// It uses `github.com/kelseyhightower/envconfig` underneath.
func FromEnviron(prefix string) Option {
	return func(c Config, m ...OptionMeta) error {
		err := expectKind(reflect.TypeOf(c), reflect.Ptr)
		if err != nil {
			return err
		}

		return env.Process(prefix, c)
	}
}
