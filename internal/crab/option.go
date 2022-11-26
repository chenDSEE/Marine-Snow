package crab

import (
	msApp "MarineSnow/pkg/app"
	msAppOption "MarineSnow/pkg/app/option"
)

type CrabOptions struct {
	Server msAppOption.ServerOption `mapstructure:"Server"`
	Log    msAppOption.LogOption    `mapstructure:"Log"`
}

var _ msApp.OptionSet = &CrabOptions{}

func NewOptions() *CrabOptions {
	// default value option can be changed in here, but crab no need to change those
	return &CrabOptions{
		Server: msAppOption.NewServerOption(),
		Log:    msAppOption.NewLogOption(),
	}
}

// NameFlagSet exposing all the flag to app framework that crab wanted to register
func (crabOpts *CrabOptions) NameFlagSet() *msApp.NameFlagSet {
	nfs := &msApp.NameFlagSet{}
	nfs.AddFlagSet(crabOpts.Server.Name(), crabOpts.Server.FlagSet())
	nfs.AddFlagSet(crabOpts.Log.Name(), crabOpts.Log.FlagSet())

	return nfs
}
