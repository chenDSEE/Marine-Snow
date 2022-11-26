package crab

import (
	msApp "MarineSnow/pkg/app"
	"fmt"
)

const cmdDescriptopion = `crab is a demo app for MarineSnow framework.
There will be a lot MarineSnow API usage example in this app.
You can find more design detail about MarineSnow in doc directory.
`

func NewApp(name string) (app *msApp.App) {
	opts := NewOptions()
	app = msApp.NewApp(name,
		msApp.WithRunFunc(runFunc(name, opts)),
		msApp.WithDescription(cmdDescriptopion),
		msApp.WithOptionSet(opts),
		msApp.WithConfigFile(""),
	)

	return app
}

func runFunc(name string, opts *CrabOptions) msApp.RunFunc {
	return func(info string) error {
		fmt.Printf("hello, I am %s. MS framework pass info[%s]\n", name, info)

		fmt.Printf("Server: IP[%s], port[%d]\n", opts.Server.IP, opts.Server.Port)
		fmt.Printf("Log: path[%s]\n", opts.Log.FullPathName)
		return nil
	}
}
