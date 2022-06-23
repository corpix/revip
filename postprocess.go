package revip

import (
	"fmt"
	"reflect"
)

func Postprocess(c Config, op ...Option) error {
	return postprocess(c, nil, op)
}

func postprocess(c Config, path []string, options []Option) error {
	value := reflect.ValueOf(c)
	valueType := reflect.TypeOf(c)
	kind := valueType.Kind()

	if kind == reflect.Ptr {
		if value.IsNil() {
			return nil // NOTE: skip nil's, this mean we don't have a default value
		}

		value = indirectValue(value)
		valueType = indirectType(valueType)
		kind = valueType.Kind()
	}

	switch kind {
	case reflect.Struct:
	case reflect.Array, reflect.Slice:
	case reflect.Map:
	default:
		return nil
	}

	//

	var err error
	for _, option := range options {
		err = option(c)
		if err != nil {
			if e, ok := err.(*ErrPostprocess); ok {
				e.Path = path
			}
			return err
		}
	}

	//

	switch kind {
	case reflect.Struct:
		return walkStruct(c, func(v reflect.Value, xs []string) error {
			return postprocess(
				v.Interface(),
				append(path, xs...),
				options,
			)
		})
	case reflect.Array, reflect.Slice:
		for n := 0; n < value.Len(); n++ {
			err := postprocess(
				value.Index(n).Interface(),
				append(path, fmt.Sprintf("[%d]", n)),
				options,
			)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, k := range value.MapKeys() {
			err := postprocess(
				value.MapIndex(k).Interface(),
				append(path, fmt.Sprintf("[%q]", k.String())),
				options,
			)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

//

func WithDefaults() Option {
	return func(c Config) error {
		var err error

		v, ok := c.(Defaultable)
		if ok && !isnil(reflect.ValueOf(v)) {
			err = expectKind(reflect.TypeOf(v), reflect.Ptr)
			if err != nil {
				return err
			}

			v.Default()
		}
		return nil
	}
}

func WithValidation() Option {
	return func(c Config) error {
		var err error

		v, ok := c.(Validatable)
		if ok && !isnil(reflect.ValueOf(v)) {
			err = expectKind(reflect.TypeOf(v), reflect.Ptr)
			if err != nil {
				return err
			}

			err = v.Validate()
			if err != nil {
				return &ErrPostprocess{
					Type: reflect.TypeOf(c).String(),
					Err:  err,
				}
			}
		}
		return nil
	}
}

func WithExpansion() Option {
	return func(c Config) error {
		var err error

		v, ok := c.(Expandable)
		if ok && !isnil(reflect.ValueOf(v)) {
			err = expectKind(reflect.TypeOf(v), reflect.Ptr)
			if err != nil {
				return err
			}

			err = v.Expand()
			if err != nil {
				return &ErrPostprocess{
					Type: reflect.TypeOf(c).String(),
					Err:  err,
				}
			}
		}
		return nil
	}
}
