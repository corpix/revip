package revip

import (
	"io/ioutil"
	"io"
	"os"
	"syscall"

	env "github.com/kelseyhightower/envconfig"
)

// Source defines generic interface for configuration source.
type Source = func(c Config) error

// FromReader is a `Source` constructor which creates a thunk
// to read configuration from `r` and decode it with `f` unmarshaler.
// Current implementation buffers all data in memory.
func FromReader(r io.Reader, f Unmarshaler) Source {
	return func(c Config) error {
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		return f(buf, c)
	}
}

// FromFile is a `Source` constructor which creates a thunk
// to read configuration from file addressable by `path` and
// decodes it with `f` unmarshaler.
func FromFile(path string, f Unmarshaler) Source {
	return func(c Config) error {
		r, err := os.Open(path)
		switch e := err.(type) {
		case *os.PathError:
			if e.Err == syscall.ENOENT {
				return &ErrFileNotFound{
					Path: path,
					Orig: err,
				}
			}
		case nil:
		default:
			return err
		}
		defer r.Close()

		return FromReader(r, f)(c)
	}
}

// FromEnviron is a `Source` constructor which creates a thunk
// to read configuration from environment.
// It uses `github.com/kelseyhightower/envconfig` underneath.
func FromEnviron(prefix string) Source {
	return func(c Config) error {
		return env.Process(prefix, c)
	}
}
