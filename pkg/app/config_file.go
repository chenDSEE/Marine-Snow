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

// WithConfigFile enable configuration feature and specified path to configuration file
func WithConfigFile(path string) OptionFunc {
	return func(app *App) {
		app.cfOption = newConfigFileOption(path)
	}
}

type configFileOption struct {
	isEnable bool
	path     string
}

// newConfigFileOption creates a new configFileOption object with default parameters and enable.
func newConfigFileOption(path string) configFileOption {
	return configFileOption{
		path:     path,
		isEnable: true,
	}
}

func (opt *configFileOption) Name() string {
	return "configFile"
}

// FlagSet create a new pflag.FlagSet for all the flag configFileOption need
func (opt *configFileOption) FlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet(opt.Name(), pflag.PanicOnError)

	fs.StringVarP(&opt.path, "config", "c", opt.path,
		"path to configuration file, support JSON, TOML, YAML, HCL, or Java properties formats.")

	return fs
}

// loadConfigFile load configuration file and replace ENV data in configuration file
func loadConfigFile(app *App, cmd *cobra.Command) error {
	viper.SetConfigFile(app.cfOption.path)

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
