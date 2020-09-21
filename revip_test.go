package revip

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRevip(t *testing.T) {
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
	c := Config{
		Foo: Foo{
			Bar: "bar",
			Qux: true,
		},
		Baz: 666,
	}

	r, err := Load(
		&c,
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
	assert.Equal(
		t,
		Config{
			Foo: Foo{Bar: "qux", Qux: false},
			Baz: 777,
			Box: []int{666, 777, 888},
		},
		c,
	)

	//

	fv := Foo{}
	err = r.Path(&fv, "Foo")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(
		t,
		Foo{Bar: "qux", Qux: false},
		fv,
	)

	//

	fvv := new(bool)
	err = r.Path(fvv, "Foo.Qux")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(
		t,
		bool(false),
		*fvv,
	)
}
