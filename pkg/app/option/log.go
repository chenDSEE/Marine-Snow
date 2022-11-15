package option

import (
	"github.com/spf13/pflag"
)

type LogOption struct {
	FullPathName string
}

// NewLogOption creates a new LogOption object with default parameters.
func NewLogOption() LogOption {
	return LogOption{
		FullPathName: "./app.log",
	}
}

func (opt *LogOption) Name() string {
	return "log"
}

func (opt *LogOption) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet(opt.Name(), pflag.PanicOnError)

	fs.StringVarP(&opt.FullPathName, "log", "l", opt.FullPathName, "specified app log file")

	return fs
}
