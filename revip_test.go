package revip

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type (
	TestConfig struct {
		Name     string              `yaml:"name"`
		Amount   int                 `yaml:"amount"`
		Provider *TestProviderConfig `yaml:"provider"`

		nonce int
	}
	TestProviderConfig struct {
		Type   string                    `yaml:"type"`
		Simple *TestSimpleProviderConfig `yaml:"simple"`
		Inline *TestInlineProviderConfig `yaml:"inline"`
	}
	TestBaseProviderConfig struct {
		Rate     int                           `yaml:"rate"`
		Auto     *bool                         `yaml:"auto"`
		Actions  []*TestActionConfig           `yaml:"actions"`
		Handlers map[string]*TestHandlerConfig `yaml:"handlers"`
	}
	TestActionConfig struct {
		Name string `yaml:"name"`
	}
	TestHandlerConfig struct {
		Name string `yaml:"name"`
	}
	TestSimpleProviderConfig struct {
		Base *TestBaseProviderConfig `yaml:"base"`
	}
	TestInlineProviderConfig struct {
		Base *TestBaseProviderConfig `yaml:",inline"`
	}
)

func (c *TestConfig) Default() {
	if c.Amount == 0 {
		c.Amount = 10
	}
	if c.Provider == nil {
		c.Provider = &TestProviderConfig{}
	}
}
func (c *TestProviderConfig) Default() {
	if c.Type == "" {
		c.Type = "simple"
	}
	if c.Type == "simple" && c.Simple == nil {
		c.Simple = &TestSimpleProviderConfig{}
	}
	if c.Type == "inline" && c.Inline == nil {
		c.Inline = &TestInlineProviderConfig{Base: &TestBaseProviderConfig{}}
	}
}
func (c *TestBaseProviderConfig) Default() {
	if len(c.Handlers) == 0 {
		c.Handlers = map[string]*TestHandlerConfig{
			"/": {Name: "root"},
		}
	}
}

func (c *TestProviderConfig) Validate() error {
	if c.Type == "" {
		return errors.New("type should not be empty")
	}
	return nil
}
func (c *TestConfig) Validate() error {
	if c.Name == "" {
		return errors.New("name should not be empty")
	}
	return nil
}

func (c *TestConfig) Expand() error {
	c.nonce += 1
	return nil
}

//

func TestConfigDefaults(t *testing.T) {
	c := &TestConfig{}
	err := Postprocess(c, WithDefaults())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "", c.Name)
	assert.Equal(t, 10, c.Amount)
	assert.NotNil(t, c.Provider)
	assert.Equal(t, "simple", c.Provider.Type)
	assert.NotNil(t, c.Provider.Simple)
	assert.Nil(t, c.Provider.Inline)
	assert.Nil(t, c.Provider.Simple.Base)
}

func TestConfigDefaultsDependant(t *testing.T) {
	c := &TestConfig{Provider: &TestProviderConfig{Type: "inline"}}
	err := Postprocess(c, WithDefaults())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "", c.Name)
	assert.Equal(t, 10, c.Amount)
	assert.NotNil(t, c.Provider)
	assert.Equal(t, "inline", c.Provider.Type)
	assert.Nil(t, c.Provider.Simple)
	assert.NotNil(t, c.Provider.Inline)
	assert.NotNil(t, c.Provider.Inline.Base)
	assert.Equal(t, "root", c.Provider.Inline.Base.Handlers["/"].Name)
}

//

func TestConfigValidation(t *testing.T) {
	c := &TestConfig{}
	err := Postprocess(c, WithValidation())
	assert.NotNil(t, err)
	assert.Equal(t, "postprocessing failed at .TestConfig: name should not be empty", err.Error())

	c.Name = "foo" // should not enter nil configurations
	err = Postprocess(c, WithValidation())
	assert.Nil(t, err)

	c.Provider = &TestProviderConfig{}
	err = Postprocess(c, WithValidation())
	assert.NotNil(t, err)
	assert.Equal(t, "postprocessing failed at .TestConfig.TestConfig.Provider: type should not be empty", err.Error())
}

//

func TestConfigNoNilPointers(t *testing.T) {
	c := &TestConfig{}
	err := Postprocess(c, WithNoNilPointers())
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, c.Provider)
	assert.NotNil(t, c.Provider.Simple)
	assert.NotNil(t, c.Provider.Simple.Base)
	assert.NotNil(t, c.Provider.Simple.Base.Actions)
	assert.NotNil(t, c.Provider.Simple.Base.Handlers)
	assert.NotNil(t, c.Provider.Inline)
	assert.NotNil(t, c.Provider.Inline.Base)
	assert.NotNil(t, c.Provider.Inline.Base.Actions)
	assert.NotNil(t, c.Provider.Inline.Base.Handlers)
}

//

func TestConfigExpansion(t *testing.T) {
	c := &TestConfig{}
	err := Postprocess(c, WithExpansion())
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, c.nonce)
}
