package app

import (
	"fmt"
	"github.com/spf13/pflag"
)

type NameFlagSet struct {
	fsMap map[string]*pflag.FlagSet
}

func (nfs *NameFlagSet) AddFlagSet(name string, fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	if nfs.fsMap == nil {
		nfs.fsMap = make(map[string]*pflag.FlagSet)
	}

	if _, existed := nfs.fsMap[name]; existed {
		msg := fmt.Sprintf("NameFlagSet.AddFlagSet(): %s FlagSet is duplicate !!!", name)
		panic(msg)
	}

	nfs.fsMap[name] = fs
}
