package crab

import (
	msApp "MarineSnow/pkg/app"
	"fmt"
)

const crabDescriptopion = `crab is a demo app for MarineSnow framework.
There will be a lot MarineSnow API usage example in this app.
You can find more design detail about MarineSnow in doc directory.
`

func NewApp(name string) *msApp.App {
	app := msApp.NewApp(name,
		msApp.WithRunFunc(runFunc(name)),
		msApp.WithDescription(crabDescriptopion),
	)

	return app
}

func runFunc(name string) msApp.RunFunc {
	return func(info string) error {
		fmt.Printf("hello, I am %s. MS framework pass info[%s]\n", name, info)
		return nil
	}
}
