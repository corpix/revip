module github.com/corpix/revip

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pelletier/go-toml v1.9.5
	github.com/stretchr/testify v1.4.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
