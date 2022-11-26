package option

import (
	"github.com/spf13/pflag"
)

type LogOption struct {
	FullPathName string `mapstructure:"path"`
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

	fs.StringVar(&opt.FullPathName, "log.path", opt.FullPathName, "specified app log file")

	return fs
}
