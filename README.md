# revip

Dead-simple configuration loader.

Usage example:

```go
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
```

Run example:

```console
user@localhost ~/p/s/g/c/revip (master)> go run ./example/main.go
config: main.Config{Foo:main.Foo{Bar:"bar", Qux:false}, Baz:666, Dox:[]string(nil), Box:[]int{666, 777, 888}}
user@localhost ~/p/s/g/c/revip (master)> REVIP_FOO_BAR=hello go run ./example/main.go
config: main.Config{Foo:main.Foo{Bar:"hello", Qux:false}, Baz:666, Dox:[]string(nil), Box:[]int{666, 777, 888}}
user@localhost ~/p/s/g/c/revip (master)> REVIP_BOX=888,777,666 go run ./example/main.go
config: main.Config{Foo:main.Foo{Bar:"bar", Qux:false}, Baz:666, Dox:[]string(nil), Box:[]int{888, 777, 666}}
```

## license

[public domain](https://unlicense.org/)
