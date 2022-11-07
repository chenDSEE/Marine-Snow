package crab

import (
	msApp "MarineSnow/pkg/app"
	"fmt"
)

func NewApp(name string) *msApp.App {
	app := msApp.NewApp(name,
		msApp.WithRunFunc(runFunc(name)),
	)

	return app
}

func runFunc(name string) msApp.RunFunc {
	return func() error {
		fmt.Printf("hello, I am %s\n", name)
		return nil
	}
}
