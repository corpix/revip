package revip

import (
	"reflect"

	"github.com/Jeffail/gabs"
	"github.com/fatih/structs"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
)

// Config is a configuration represented by user-specified type.
type Config = interface{}

// Revip represents loaded by `Load` configuration.
type Revip struct {
	// config represents configuration data, it should always be a pointer.
	config Config
}

// Unwrap returns a pointer to the inner configuration data structure.
func (r *Revip) Unwrap() interface{} { return r.config }

// Config writes a shallow copy of the configuration into `v`.
func (r *Revip) Config(v interface{}) error {
	return copier.Copy(v, r.config)
}

// Path uses `github.com/Jeffail/gabs` to retrieve configuration key
// or sub-tree into `v` which is addressable by provided `path` or
// return an error if key was not found(`ErrNotFound`) or
// something gone terribly wrong.
func (r *Revip) Path(v Config, path string) error {
	g, err := gabs.Consume(structs.Map(r.config))
	if err != nil {
		return err
	}

	if !g.ExistsP(path) {
		return &ErrPathNotFound{Path: path}
	}

	p := g.Path(path).Data()

	err = mapstructure.WeakDecode(p, v)
	if err != nil {
		return err
	}

	return nil
}

// New wraps configuration represented by `c` with come useful methods.
func New(c Config) *Revip {
	if reflect.TypeOf(c).Kind() != reflect.Ptr {
		panic("config must be a pointer")
	}

	return &Revip{config: c}
}
