package revip

import (
	"testing"
	"fmt"
	"bytes"

	//"github.com/stretchr/testify/assert"
)

func TestRevip(t *testing.T) {
	// bash -c 'REVIP_BAZ=777 REVIP_FOO_BAR=qux REVIP_DOX=1,2,3 go test -v ./...'

	type Foo struct {
		Bar string
		Qux bool
	}
	type Config struct {
		Foo Foo
		Baz int
		Dox []string
		Box []int
	}
	config := Config{
		Foo: Foo{
			Bar: "bar",
			Qux: true,
		},
		Baz: 666,
	}

	err := Unmarshal(
		&config,
		FromReader(
			bytes.NewBuffer([]byte(`foo: { qux: false }`)),
			YamlUnmarshaler,
		),
		FromReader(
			bytes.NewBuffer([]byte(`box = [666,777,888]`)),
			TomlUnmarshaler,
		),
		FromEnviron("revip"),
	)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("-- %#v\n", config)
	// -- revip.Config{Foo:revip.Foo{Bar:"qux", Qux:false}, Baz:777, Dox:[]string{"1", "2", "3"}, Box:[]int{666, 777, 888}}
}
