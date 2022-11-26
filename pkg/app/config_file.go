// Configuration file feature is enable by default. And the path to configuration
// is specified by '-c /PATH/TO/FILE/FILE_NAME.FILE_TYPE'
// Disable this by WithNoConfigFile().
// configFile is an option, but not exposed, only used by app framework.

package app

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// WithNoConfigFile disable configuration file feature
func WithNoConfigFile() OptionFunc {
	return func(app *App) {
		app.noConfigFile = true
	}
}

const (
	defaultConfigFile = "./config.yaml"
)

type configFileOption struct {
	name string
}

// newConfigFileOption creates a new ServerOption object with default parameters.
func newConfigFileOption() configFileOption {
	return configFileOption{
		name: defaultConfigFile,
	}
}

func (opt *configFileOption) Name() string {
	return "configFile"
}

// FlagSet create a new pflag.FlagSet for all the flag configFileOption need
func (opt *configFileOption) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet(opt.Name(), pflag.PanicOnError)

	fs.StringVarP(&opt.name, "config", "c", opt.name,
		"path to configuration file, support JSON, TOML, YAML, HCL, or Java properties formats.")

	return fs
}

func loadConfigFile(app *App, cmd *cobra.Command) error {
	viper.SetConfigFile(app.cfOption.name)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// update flag
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	if err := viper.Unmarshal(app.optionSet); err != nil {
		return err
	}

	return nil
}
