package crab

import (
	msApp "MarineSnow/pkg/app"
	msAppOption "MarineSnow/pkg/app/option"
)

type crabOptions struct {
	server msAppOption.ServerOption
	log    msAppOption.LogOption
}

var _ msApp.OptionSet = &crabOptions{}

func NewOptions() msApp.OptionSet {
	// default value option can be changed in here, but crab no need to change those
	return &crabOptions{
		server: msAppOption.NewServerOption(),
		log:    msAppOption.NewLogOption(),
	}
}

// NameFlagSet exposing all the flag to app framework that crab wanted to register
func (crabOpts *crabOptions) NameFlagSet() *msApp.NameFlagSet {
	nfs := &msApp.NameFlagSet{}
	nfs.AddFlagSet(crabOpts.server.Name(), crabOpts.server.FlagSet())
	nfs.AddFlagSet(crabOpts.log.Name(), crabOpts.log.FlagSet())

	return nfs
}
