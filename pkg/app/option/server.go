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

func (opt *ServerOption) Name() string {
	return "server"
}

func (opt *ServerOption) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet(opt.Name(), pflag.PanicOnError)

	fs.IPVarP(&opt.IP, "IP", "i", net.ParseIP(defaultHost), "server host listen to")
	fs.IntVarP(&opt.Port, "port", "p", 80, "server port listen to")

	return fs
}
