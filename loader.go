package revip

import (
	"reflect"

	"github.com/fatih/structs"

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

func Postprocess(m Config) error { return postprocess(m, nil) }

func postprocessApply(m Config, path []string) error {
	de, ok := m.(Defaultable)
	if ok {
		de.Default()
	}

	ve, ok := m.(Validatable)
	if ok {
		err := ve.Validate()
		if err != nil {
			return &ErrPostprocess{
				Type: reflect.TypeOf(m).String(),
				Path: path,
				Err:  err,
			}
		}
	}
	return nil
}

func postprocess(m Config, path []string) error {
	err := postprocessApply(m, path)
	if err != nil {
		return err
	}

	//

	t := reflect.TypeOf(m)

	if indirectType(t).Kind() != reflect.Struct {
		return nil
	}

	//

	for _, v := range structs.Fields(m) {
		if !v.IsExported() {
			continue
		}

		err := postprocess(
			v.Value(),
			append(path, v.Name()),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
