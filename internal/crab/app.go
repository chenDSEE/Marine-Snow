package crab

import (
	msApp "MarineSnow/pkg/app"
	"errors"
	"fmt"
	"github.com/spf13/viper"
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
		//msApp.WithNoConfigFile(),
	)

	return app
}

func runFunc(name string, optSet msApp.OptionSet) msApp.RunFunc {
	return func(info string) error {
		fmt.Printf("hello, I am %s. MS framework pass info[%s]\n", name, info)
		crabOpts, ok := optSet.(*CrabOptions)
		if !ok {
			return errors.New("error OptionSet")
		}
		fmt.Printf("2 ----> Server port:%d\n", viper.GetInt("Server.port"))

		fmt.Printf("Server: IP[%s], port[%d]\n", crabOpts.Server.IP, crabOpts.Server.Port)
		fmt.Printf("Log: path[%s]\n", crabOpts.Log.FullPathName)
		return nil
	}
}
