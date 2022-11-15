package app

import (
	"github.com/spf13/pflag"

	msAppOption "MarineSnow/pkg/app/option"
)

// Option is a feature that app instance can enable or not.
type Option interface {
	Name() string
	FlagSet() *pflag.FlagSet
}

var _ Option = &msAppOption.ServerOption{}
var _ Option = &msAppOption.LogOption{}

// OptionSet is an abstraction for app framework to access feature that app instance wanting to enable.
// this interface should be implemented by all app instance needing configuration file or Command flag
type OptionSet interface {
	// NameFlagSet return all the pflag.FlagSet need by app instance,
	// app abstract framework will register them into cobra Command
	NameFlagSet() *NameFlagSet
}

// WithOptionSet to enable the feature that app instance wanted to enable
func WithOptionSet(optSet OptionSet) OptionFunc {
	return func(app *App) {
		app.optionSet = optSet

		// if OptionSet.NameFlagSet() create and return a new NameFlagSet object.
		// OptionSet.NameFlagSet() may make a side effect to the flag variables
		// so this method can only call once in app framework.
		app.nfs = app.optionSet.NameFlagSet()
	}
}
