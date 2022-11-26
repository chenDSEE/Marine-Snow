package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RunFunc Q&A:
// should pkg/app which play as an app framework, return any framework level struct backup to app ?
// RunFunc should contain any app framework level struct as param ?
type RunFunc func(string) error
type OptionFunc func(*App)

type App struct {
	/* APP information */
	Name        string
	description string

	/* APP command */
	runFunc RunFunc
	rootCmd *cobra.Command

	/* APP option and flag */
	optionSet    OptionSet
	nfs          *NameFlagSet
	noConfigFile bool
	cfOption     configFileOption
}

func WithRunFunc(fun RunFunc) OptionFunc {
	return func(app *App) {
		app.runFunc = fun
	}
}

func WithDescription(desc string) OptionFunc {
	return func(app *App) {
		app.description = desc
	}
}

func NewApp(name string, opts ...OptionFunc) *App {
	app := &App{
		Name: name,
	}

	/* setup option for app */
	for _, fun := range opts {
		fun(app)
	}

	/* setup command for app */
	app.buildCommand()

	// app framework option
	if app.noConfigFile == false {
		app.cfOption = newConfigFileOption()
		app.rootCmd.PersistentFlags().AddFlagSet(app.cfOption.FlagSet())
		//app.nfs.AddFlagSet(app.cfOption.Name(), app.cfOption.FlagSet())
	}

	return app
}

func (app *App) buildCommand() {
	/* build basic root command */
	app.rootCmd = &cobra.Command{
		Use:   app.Name,
		Short: "application " + app.Name,
		Long:  app.description,
	}

	for _, fs := range app.nfs.fsMap {
		app.rootCmd.Flags().AddFlagSet(fs)
	}

	// Only run app with App.RunFunc register
	// Or return by cobra.Command.Execute()
	if app.runFunc != nil {
		app.rootCmd.RunE = app.runCommand
	}
}

func (app *App) Run() {
	if err := app.rootCmd.Execute(); err != nil {
		fmt.Printf("%s fail to run, exit with: %s\n", app.Name, err.Error())
		os.Exit(1)
	}
}

func (app *App) runCommand(cmd *cobra.Command, args []string) error {
	/* load configuration data from configuration file */
	if app.noConfigFile == false && len(app.cfOption.name) != 0 {
		if err := loadConfigFile(app, cmd); err != nil {
			return err
		}
	}

	/* parse command, flag and configuration */
	DebugAppDump(app)

	/* execute app registered RunFunc */
	if app.runFunc != nil {
		return app.runFunc("run success")
	}

	return errors.New("exit command run without App.RunFunc execute")
}
