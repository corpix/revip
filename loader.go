package revip

import (
	json "encoding/json"
	yaml "github.com/go-yaml/yaml"
	toml "github.com/pelletier/go-toml"
)

// Unmarshaler describes a generic unmarshal interface
// which could be used to extend supported formats by defining new `Option`
// implementations.
type Unmarshaler = func(in []byte, v interface{}) error

var (
	JsonUnmarshaler Unmarshaler = json.Unmarshal
	YamlUnmarshaler Unmarshaler = yaml.Unmarshal
	TomlUnmarshaler Unmarshaler = toml.Unmarshal
)

// Load applies each `op` in order to fill the configuration in `v` and
// constructs a `*Revip` data-structure.
func Load(v Config, op ...Source) (*Revip, error) {
	var err error
	for _, f := range op {
		err = f(v)
		if err != nil {
			return nil, err
		}
	}

	return New(v), nil
}
