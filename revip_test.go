package revip

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FooSimple struct {
	Bar string
	Qux bool
}
type Untouched struct{}
type ConfigSimple struct {
	Foo FooSimple
	Baz int
	Dox []string
	Box []int
}

func TestRevipSimple(t *testing.T) {
	c := ConfigSimple{
		Foo: FooSimple{
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
		ConfigSimple{
			Foo: FooSimple{Bar: "qux", Qux: false},
			Baz: 777,
			Box: []int{666, 777, 888},
		},
		c,
	)

	//

	fv := FooSimple{}
	err = r.Path(&fv, "Foo")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(
		t,
		FooSimple{Bar: "qux", Qux: false},
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

//

var (
	fooPostprocessDefaultCalled  = 0
	fooPostprocessValidateCalled = 0
)

type FooPostprocess struct {
	Bar string
	Qux bool
}

func (f *FooPostprocess) Default() {
	fooPostprocessDefaultCalled++
}
func (f *FooPostprocess) Validate() error {
	fooPostprocessValidateCalled++
	return nil
}

var (
	configPostprocessDefaultCalled  = 0
	configPostprocessValidateCalled = 0
)

type ConfigPostprocess struct {
	Foo *FooPostprocess
	Baz int
	Dox []string
	Box []int
}

func (f *ConfigPostprocess) Default() {
	configPostprocessDefaultCalled++
}
func (f *ConfigPostprocess) Validate() error {
	configPostprocessValidateCalled++
	return nil
}

func TestRevipPostprocess(t *testing.T) {
	c := ConfigPostprocess{
		Foo: &FooPostprocess{
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

	err = Postprocess(
		&c,
		WithDefaults(),
		WithValidation(),
	)
	if err != nil {
		t.Error(err)
	}

	//

	assert.Equal(t, 1, configPostprocessDefaultCalled)
	assert.Equal(t, 1, configPostprocessValidateCalled)
	assert.Equal(t, 1, fooPostprocessDefaultCalled)
	assert.Equal(t, 1, fooPostprocessValidateCalled)

	//

	assert.Equal(
		t,
		ConfigPostprocess{
			Foo: &FooPostprocess{Bar: "qux", Qux: false},
			Baz: 777,
			Box: []int{666, 777, 888},
		},
		c,
	)

	//

	fv := FooPostprocess{}
	err = r.Path(&fv, "Foo")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(
		t,
		FooPostprocess{Bar: "qux", Qux: false},
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
