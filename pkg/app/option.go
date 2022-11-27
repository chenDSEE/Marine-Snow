package app

import (
	"github.com/spf13/pflag"

	msAppOption "MarineSnow/pkg/app/option"
)

// Option is a feature that app instance can enable or not.
type Option interface {
	Name() string
	FlagSet() *pflag.FlagSet // all the flag set of this Option, binding to the variables
	Validate() []error       // Validate checks validation of this Option
}

var _ Option = &msAppOption.ServerOption{}
var _ Option = &msAppOption.LogOption{}

var _ Option = &configFileOption{} // not exposed

// OptionSet is an abstraction for app framework to access feature that app instance wanting to enable.
// this interface should be implemented by all app instance needing configuration file or Command flag
type OptionSet interface {
	// NameFlagSet return all the pflag.FlagSet need by app instance,
	// app abstract framework will register them into cobra Command
	NameFlagSet() *NameFlagSet
	Validate() []error // Validate checks all the Option validation in this OptionSet
}

// WithOptionSet to enable the feature that app instance wanted to enable
// OptionSet must be provided by user, or user is hard to access the related configuration data.
func WithOptionSet(optSet OptionSet) OptionFunc {
	return func(app *App) {
		app.optionSet = optSet

		// if OptionSet.NameFlagSet() create and return a new NameFlagSet object.
		// OptionSet.NameFlagSet() may make a side effect to the flag variables
		// so this method can only call once in app framework.
		app.nfs = app.optionSet.NameFlagSet()
	}
}
