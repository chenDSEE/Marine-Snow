package app

import (
	"fmt"
	"os"
)

// RunFunc Q&A:
// should pkg/app which play as an app framework, return any framework level struct backup to app ?
// RunFunc should contain any app framework level struct as param ?
type RunFunc func() error
type OptionFunc func(*App)

type App struct {
	Name    string
	runFunc RunFunc
}

func NewApp(name string, opts ...OptionFunc) *App {
	app := &App{
		Name: name,
	}

	/* setup option for app */
	for _, fun := range opts {
		fun(app)
	}

	return app
}

func WithRunFunc(fun RunFunc) OptionFunc {
	return func(app *App) {
		app.runFunc = fun
	}
}

func (app *App) Run() {
	if app.runFunc != nil {
		if err := app.runFunc(); err != nil {
			fmt.Printf("App exit with error:[%s]\n", err.Error())
			os.Exit(1)
		}
	}
}
