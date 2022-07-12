package revip

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreePathString(t *testing.T) {
	type (
		Bar struct {
			D int
			E string
		}
		Foo struct {
			A map[int]Bar
			B map[string]Bar
			C []Bar
			D Bar
			E *Bar
		}
	)

	var paths []string

	_, _ = NewTree(
		reflect.ValueOf(Foo{
			A: map[int]Bar{1: {}},
			B: map[string]Bar{"hello": {}, "world": {}},
			C: []Bar{{}, {}},
			E: &Bar{},
		}),
		func(tt Tree) error {
			paths = append(paths, TreePathString(tt))
			return nil
		},
	)

	assert.Equal(t, []string{
		".Foo",
		".Foo.A",
		".Foo.A[1]",
		".Foo.A[1].D",
		".Foo.A[1].E",
		".Foo.B",
		".Foo.B[hello]",
		".Foo.B[hello].D",
		".Foo.B[hello].E",
		".Foo.B[world]",
		".Foo.B[world].D",
		".Foo.B[world].E",
		".Foo.C",
		".Foo.C[0]",
		".Foo.C[0].D",
		".Foo.C[0].E",
		".Foo.C[1]",
		".Foo.C[1].D",
		".Foo.C[1].E",
		".Foo.D",
		".Foo.D.D",
		".Foo.D.E",
		".Foo.E",
		".Foo.E.Bar",
		".Foo.E.Bar.D",
		".Foo.E.Bar.E",
	}, paths)
}
