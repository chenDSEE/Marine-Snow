package config

import (
	"fmt"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
	"time"
)

// NOTE: There is little difference to access configuration in code: GetConfigurationString("aaa.bbb.ccc")
// or ConfigObj.AAA.BBB.CCC
//
// when you change the configuration struct, you must change code like:
//    GetConfigurationStringSlice("aaa.bbb.ccc") to  GetConfigurationStringSlice("aaa.bbb.ccc")
// or GetConfigurationString("aaa.bbb.ccc")      to  GetConfigurationString("aaa.bbb.sss.ccc")
// or ConfigObj.AAA.BBB.CCC                      to  configObj.AAA.BBB.SSS.CCC
//
// But via GetConfigurationString("aaa.bbb.ccc") have to decode the string everytime you call.
// To access configuration via ConfigObj.AAA.BBB.CCC will be better here. For convenience you can define a global
// function as shortcut for specified configuration. For example: GetAB_CCC() stand for ConfigObj.aaa.bbb.ccc

type yamlParser struct {
	cfgType string
}

func (p yamlParser) Unmarshal(in []byte, out interface{}) error {
	// TODO: update Placeholder with OS environment variable

	err := yaml.Unmarshal(in, out)
	if err != nil {
		return err
	}

	return nil
}

func (p yamlParser) SupportFileType() []string {
	return []string{"yaml", "yml"}
}

type Decoder struct {
	path    string
	file    string
	cfgType string
	parser  yamlParser

	checkInterval   time.Time
	updateChan      chan struct{}
	isEnableMonitor bool
}

// TODO: Divide configuration into app part and framework part will be better
func NewDecoder(where string, cfgType string) *Decoder {
	decoder := &Decoder{
		cfgType: cfgType,
		parser:  yamlParser{},
	}

	isFile := false
	part := strings.Split(where, ".")
	if len(part) >= 1 {
		for _, fileType := range decoder.parser.SupportFileType() {
			if part[len(part)-1] == fileType {
				isFile = true
				break
			}
		}
	}

	if isFile {
		decoder.file = where
		decoder.path = ""
	} else {
		decoder.file = ""
		decoder.path = where
	}

	return decoder
}

func (decoder *Decoder) MonitorEnable() {
	panic("impl me")
}

func (decoder *Decoder) LoadConfig(out interface{}) error {
	return loadConfigFile(decoder, decoder.file, out)
}

func loadConfigFile(decoder *Decoder, pathTofile string, out interface{}) error {
	fmt.Printf("file:[%s]\n", pathTofile)
	buf, err := ioutil.ReadFile(pathTofile)
	if err != nil {
		return err
	}

	err = decoder.parser.Unmarshal(buf, out)
	if err != nil {
		return err
	}

	return nil
}

// TODO: pass the parsed result out via channel(with reflact)
func (decoder *Decoder) UpdateConfig() <-chan struct{} {
	panic("impl me")
	return nil
}
