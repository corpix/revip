package revip

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/kelseyhightower/envconfig"
	toml "github.com/pelletier/go-toml"
)

//

type (
	//Marshaler = func(v interface{}) ([]byte, error)
	Unmarshaler = func(in []byte, v interface{}) error
)

var (
	JsonUnmarshaler Unmarshaler = json.Unmarshal
	YamlUnmarshaler Unmarshaler = yaml.Unmarshal
	TomlUnmarshaler Unmarshaler = toml.Unmarshal
)

//

type Option = func(c *Revip) error

func FromReader(r io.Reader, f Unmarshaler) func(*Revip) error {
	return func(c *Revip) error {
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		return f(buf, c.config)
	}
}

func FromFile(path string, f Unmarshaler) func(*Revip) error {
	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return FromReader(r, f)
}

func FromEnviron(prefix string) func(*Revip) error {
	return func(c *Revip) error {
		return envconfig.Process(prefix, c.config)
	}
}

//

type Revip struct {
	config interface{}
}

//

func Unmarshal(v interface{}, op ...Option) error {
	var (
		r   = &Revip{config: v}
		err error
	)
	for _, f := range op {
		err = f(r)
		if err != nil {
			return err
		}
	}

	return nil
}
