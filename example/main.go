package main

import (
	"fmt"
	"bytes"

	"github.com/corpix/revip"
)

type (
	Foo struct {
		Bar string
		Qux bool
	}
	Config struct {
		Foo Foo
		Baz int
		Dox []string
		Box []int
	}
)

func main() {
	c := Config{
		Foo: Foo{
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

	fmt.Printf("config: %#v\n", c)
}
