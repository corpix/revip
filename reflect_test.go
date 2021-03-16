package revip

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestWalkStruct(t *testing.T) {
	samples := []struct {
		s     interface{}
		kinds []reflect.Kind
		paths [][]string
		res   func([]string) error
		err   error
	}{
		{
			s: struct {
				Foo       string
				Bar       int
				invisible int
			}{},
			kinds: []reflect.Kind{reflect.String, reflect.Int},
			paths: [][]string{{"Foo"}, {"Bar"}},
		},
		{
			s: struct {
				Foo string
				Bar *struct{ Baz string }
			}{},
			kinds: []reflect.Kind{reflect.String, reflect.Ptr},
			paths: [][]string{{"Foo"}, {"Bar"}},
		},
		{
			s: struct {
				Foo string
				Bar *struct{ Baz string }
			}{
				Bar: &struct{ Baz string }{},
			},
			kinds: []reflect.Kind{reflect.String, reflect.Ptr, reflect.Struct, reflect.String},
			paths: [][]string{{"Foo"}, {"Bar"}, {"Bar"}, {"Bar", "Baz"}}, // bar duplicated because we visit pointer first
		},
		{
			s: struct {
				Foo string
				Bar struct{ Baz struct{ Qux int } }
				Dox int
			}{},
			kinds: []reflect.Kind{reflect.String, reflect.Struct, reflect.Struct, reflect.Int},
			paths: [][]string{{"Foo"}, {"Bar"}, {"Bar", "Baz"}, {"Dox"}},
			res: func(path []string) error {
				if "Bar.Baz" == strings.Join(path, ".") {
					return skipBranch
				}
				return nil
			},
		},
		{
			s: struct {
				Foo string
				Bar int
			}{},
			kinds: []reflect.Kind{reflect.String},
			paths: [][]string{{"Foo"}},
			res: func(path []string) error {
				if path[0] == "Foo" {
					return stopIteration
				}
				return nil
			},
		},
	}
	for k, sample := range samples {
		t.Run(fmt.Sprintf("%d", k), func(t *testing.T) {
			kinds := []reflect.Kind{}
			paths := [][]string{}

			iter := func(v reflect.Value, path []string) error {
				kinds = append(kinds, v.Kind())
				paths = append(paths, path)

				if sample.res != nil {
					return sample.res(path)
				}
				return nil
			}

			msg := spew.Sdump(sample)
			assert.Equal(t, sample.err, walkStruct(sample.s, iter), msg)
			assert.Equal(t, sample.kinds, kinds, msg)
			assert.Equal(t, sample.paths, paths, msg)
		})
	}
}
