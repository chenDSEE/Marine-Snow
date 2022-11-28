package option

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
)

const (
	defaultHost = "127.0.0.1"
)

type ServerOption struct {
	IP   string `mapstructure:"IP"`
	Port int    `mapstructure:"Port"`
}

// NewServerOption creates a new ServerOption object with default parameters.
func NewServerOption() ServerOption {
	return ServerOption{
		IP:   defaultHost,
		Port: 80,
	}
}

func (opt *ServerOption) Name() string {
	return "server"
}

// FlagSet create a new pflag.FlagSet for all the flag ServerOption need.
func (opt *ServerOption) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet(opt.Name(), pflag.PanicOnError)

	fs.StringVar(&opt.IP, "server.ip", opt.IP, "server host listen to")
	fs.IntVar(&opt.Port, "server.port", opt.Port, "server port listen to")

	return fs
}

func (opt *ServerOption) Validate() []error {
	var errs []error

	if opt.Port < 0 {
		msg := fmt.Sprintf("ServerOption: port[%d], this should be positive", opt.Port)
		errs = append(errs, errors.New(msg))
	}

	return errs
}
