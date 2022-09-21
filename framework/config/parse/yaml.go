package parse

import (
	yaml "gopkg.in/yaml.v3"
)

type yamlParser struct {
	cfgType string
}

func (p yamlParser) Unmarshal(in []byte, out interface{}) error {
	return yaml.Unmarshal(in, out)
}

func (p yamlParser) SupportFileType() []string {
	return []string{"yaml", "yml"}
}
