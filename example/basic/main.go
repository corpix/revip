package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/corpix/revip"

	"github.com/davecgh/go-spew/spew"
)

type Config struct {
	// yaml keys must be all lower-case
	// otherwise you need to tag every field
	// see: https://github.com/go-yaml/yaml/issues/123
	SerialNumber int `yaml:"serialNumber"`

	Nested      *NestedConfig
	MapNested   map[string]*NestedConfig
	SliceNested []*NestedConfig

	StringSlice []string
	IntSlice    []int

	key string
}

func (c *Config) Default() {
loop:
	for {
		switch {
		case c.Nested == nil:
			c.Nested = &NestedConfig{}
		case c.MapNested == nil:
			c.MapNested = map[string]*NestedConfig{}
		case c.SliceNested == nil:
			c.SliceNested = []*NestedConfig{}
		default:
			break loop
		}
	}
}

func (c *Config) Validate() error {
	if c.Nested.Flag {
		return fmt.Errorf("nested flag should be false")
	}
	if c.SerialNumber <= 0 {
		return fmt.Errorf("serialNumber should be greater than zero")
	}
	if len(c.IntSlice) != 3 {
		return fmt.Errorf("intSlice length should be 3")
	}
	return nil
}

func (c *Config) Expand() error {
	buf, err := ioutil.ReadFile("./key")
	if err != nil {
		return err
	}
	c.key = string(buf)

	return nil
}

//

type NestedConfig struct {
	Value string
	Flag  bool
}

func (c *NestedConfig) Default() {
loop:
	for {
		switch {
		case c.Value == "":
			c.Value = "default"
		default:
			break loop
		}
	}
}

//

func main() {
	c := Config{
		Nested: &NestedConfig{
			Value: "hello world",
			Flag:  true,
		},
	}

	_, err := revip.Load(
		&c,
		revip.FromReader(
			bytes.NewBuffer([]byte(`{"nested":{"flag": false}}`)),
			revip.JsonUnmarshaler,
		),
		revip.FromReader(
			bytes.NewBuffer([]byte(`serialNumber: 1`)),
			revip.YamlUnmarshaler,
		),
		revip.FromReader(
			bytes.NewBuffer([]byte(`intSlice = [666,777,888]`)),
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
		revip.WithExpansion(),
	)
	if err != nil {
		panic(err)
	}

	spew.Dump(c)
}
