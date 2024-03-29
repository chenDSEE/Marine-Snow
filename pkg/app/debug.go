package app

import (
	"fmt"
	"github.com/spf13/pflag"
)

// DebugAppDump just for debug purpose to dump App.
func DebugAppDump(app *App) {
	fmt.Printf("\n====== app.App dump ======\n")
	fmt.Printf("Name: %s\n", app.Name)
	fmt.Printf("description: %s\n", app.description)
	fmt.Printf("configuration file: isEnable[%v], path[%s]\n", app.cfOption.isEnable, app.cfOption.path)
	fmt.Printf("OptionSet:\n")
	DebugNameFlagSetDump(app.nfs)
	fmt.Printf("==========================\n\n")
}

func DebugNameFlagSetDump(nfs *NameFlagSet) {
	for name, fs := range nfs.fsMap {
		fmt.Printf("  [%s] option:\n", name)
		fs.VisitAll(func(flag *pflag.Flag) {
			fmt.Printf("    --%s: %s\n", flag.Name, flag.Value.String())
		})
		fmt.Printf("\n")
	}
}
