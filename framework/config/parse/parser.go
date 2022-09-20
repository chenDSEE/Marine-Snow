package parse

type CfgParser interface {
	Unmarshal(in []byte, out interface{}) error
	SupportFileType() []string
}

func NewCfgParser(cfgType string) CfgParser {
	switch cfgType {
	case "yaml":
		return yamlParser{cfgType: cfgType}
	default:
		return nil
	}
}
