package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/corpix/revip"

	"github.com/davecgh/go-spew/spew"
)

type (
	Config struct {
		// yaml keys must be all lower-case
		// otherwise you need to tag every field
		// see: https://github.com/go-yaml/yaml/issues/123
		SerialNumber int `yaml:"serialNumber"`

		Nested      *NestedConfig
		MapNested   map[string]*NestedConfig
		SliceNested []*NestedConfig

		StringSlice []string
		IntSlice    []int

		*EmbeddedConfig `yaml:",inline,omitempty"`

		key string
	}
	NestedConfig struct {
		Value string
		Flag  bool
	}
	EmbeddedConfig struct {
		EmbeddedStrField      string `yaml:"str"`
		EmbeddedIntField      int    `yaml:"int"`
		*DeeplyEmbeddedConfig `yaml:",inline,omitempty"`
	}
	DeeplyEmbeddedConfig struct {
		DeeplyEmbeddedStrField string `yaml:"deep-str"`
	}
)

func (c *Config) Default() {
	if c.Nested == nil {
		c.Nested = &NestedConfig{}
	}
	if c.MapNested == nil {
		c.MapNested = map[string]*NestedConfig{}
	}
	if c.SliceNested == nil {
		c.SliceNested = []*NestedConfig{}
	}
	if c.EmbeddedConfig == nil {
		c.EmbeddedConfig = &EmbeddedConfig{}
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

func (c *EmbeddedConfig) Default() {
	if c.EmbeddedStrField == "" {
		c.EmbeddedStrField = "embedded field"
	}
}
func (c *EmbeddedConfig) Validate() error {
	// shadow Validate() inherited from DeeplyEmbeddedConfig
	// or validation will fail, because default value set in Default()
	// method and we Postprocess() tree level by level from top to bottom
	return nil
}

func (c *DeeplyEmbeddedConfig) Default() {
	if c.DeeplyEmbeddedStrField == "" {
		c.DeeplyEmbeddedStrField = "deeply embedded field"
	}
}

func (c *DeeplyEmbeddedConfig) Validate() error {
	if c.DeeplyEmbeddedStrField == "" {
		return fmt.Errorf("deep-str should not be empty")
	}
	return nil
}

//

func main() {
	c := &Config{
		Nested: &NestedConfig{
			Value: "hello world",
			Flag:  true,
		},
	}

	_, err := revip.Load(
		c,
		revip.FromReader(
			bytes.NewBuffer([]byte(`{"nested":{"flag": false}}`)),
			revip.JsonUnmarshaler,
		),
		revip.FromReader(
			bytes.NewBuffer([]byte(`{serialNumber: 1, int: 666}`)),
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
		c,
		revip.WithDefaults(),
		revip.WithValidation(),
		revip.WithExpansion(),
	)
	if err != nil {
		panic(err)
	}

	spew.Dump(c)
}
