package option

import (
	"github.com/spf13/pflag"
)

type LogOption struct {
	FullPathName string
}

func (opt *LogOption) Name() string {
	return "log"
}

func (opt *LogOption) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet(opt.Name(), pflag.PanicOnError)

	fs.StringVarP(&opt.FullPathName, "log", "l", "./log.log", "specified app log file")

	return fs
}
