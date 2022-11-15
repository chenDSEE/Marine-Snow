package option

import (
	"github.com/spf13/pflag"
	"net"
)

const (
	defaultHost = "127.0.0.1"
)

type ServerOption struct {
	IP   net.IP
	Port int
}

// NewServerOption creates a new ServerOption object with default parameters.
func NewServerOption() ServerOption {
	return ServerOption{
		IP:   net.ParseIP(defaultHost),
		Port: 80,
	}
}

func (opt *ServerOption) Name() string {
	return "server"
}

// FlagSet create a new pflag.FlagSet for all the flag ServerOption need
func (opt *ServerOption) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet(opt.Name(), pflag.PanicOnError)

	fs.IPVarP(&opt.IP, "IP", "i", opt.IP, "server host listen to")
	fs.IntVarP(&opt.Port, "port", "p", opt.Port, "server port listen to")

	return fs
}
