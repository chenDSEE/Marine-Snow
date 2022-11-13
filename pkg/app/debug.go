package app

import (
	"fmt"
	"github.com/spf13/pflag"
)

// DebugAppDump just for debug purpose to dump App
func DebugAppDump(app *App) {
	fmt.Printf("\n====== app.App dump ======\n")
	fmt.Printf("Name: %s\n", app.Name)
	fmt.Printf("description: %s\n", app.description)
	fmt.Printf("OptionSet:\n")
	DebugOptionSetDump(app.optionSet)
	fmt.Printf("==========================\n\n")
}

func DebugOptionSetDump(optSet OptionSet) {
	for name, fs := range optSet.NameFlagSet().fsMap {
		fmt.Printf("  [%s] option:\n", name)
		fs.VisitAll(func(flag *pflag.Flag) {
			fmt.Printf("    --%s: %s\n", flag.Name, flag.Value.String())
		})
		fmt.Printf("\n")
	}
}
