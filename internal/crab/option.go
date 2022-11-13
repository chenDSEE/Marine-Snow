package crab

import (
	msApp "MarineSnow/pkg/app"
	msAppOption "MarineSnow/pkg/app/option"
)

type crabOptions struct {
	nfs    *msApp.NameFlagSet
	server msAppOption.ServerOption
	log    msAppOption.LogOption
}

var _ msApp.OptionSet = &crabOptions{}

func NewOptions() msApp.OptionSet {
	return &crabOptions{}
}

// NameFlagSet exposing all the flag to app framework that crab wanted to register
func (crabOpts *crabOptions) NameFlagSet() *msApp.NameFlagSet {
	if crabOpts.nfs == nil {
		crabOpts.nfs = &msApp.NameFlagSet{}
		crabOpts.nfs.AddFlagSet(crabOpts.server.Name(), crabOpts.server.FlagSet())
		crabOpts.nfs.AddFlagSet(crabOpts.log.Name(), crabOpts.log.FlagSet())
	}

	return crabOpts.nfs
}
