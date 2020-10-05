package main

import (
	"bytes"
	"fmt"

	"github.com/corpix/revip"
)

type (
	Foo struct {
		Bar string
		Qux bool
	}
	Config struct {
		Foo *Foo
		Baz int
		Dox []string
		Box []int
	}
)

func (c *Config) Validate() error {
	if c.Baz <= 0 {
		return fmt.Errorf("baz should be greater than zero")
	}
	return nil
}

func (c *Config) Default() {
loop:
	for {
		switch {
		case c.Foo == nil:
			c.Foo = &Foo{Bar:"bar default", Qux: true}
		default:
			break loop
		}
	}
}

func main() {
	c := Config{
		Foo: &Foo{
			Bar: "bar",
			Qux: true,
		},
		Baz: 666,
	}

	_, err := revip.Load(
		&c,
		revip.FromReader(
			bytes.NewBuffer([]byte(`{"foo":{"qux": false}}`)),
			revip.JsonUnmarshaler,
		),
		revip.FromReader(
			bytes.NewBuffer([]byte(`box = [666,777,888]`)),
			revip.TomlUnmarshaler,
		),
		revip.FromEnviron("revip"),
	)
	if err != nil {
		panic(err)
	}

	err = revip.Postprocess(
		&c,
		revip.WithDefaults(),
		revip.WithValidation(),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("config: %#v\n", c)
}
