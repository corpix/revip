module github.com/corpix/revip

go 1.17

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pelletier/go-toml v1.9.5
	github.com/stretchr/testify v1.7.1
)

require (
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace gopkg.in/yaml.v2 v2.4.0 => github.com/corpix/yaml v0.0.0-20220706142132-cc85bc9b2427
