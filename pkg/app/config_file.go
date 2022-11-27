// Configuration file feature is enable by default. And the path to configuration
// is specified by '-c /PATH/TO/FILE/FILE_NAME.FILE_TYPE'
// Disable this by WithNoConfigFile().
// configFile is an option, but not exposed, only used by app framework.

package app

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
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

func (opt *configFileOption) Validate() []error {
	// no need to check, do nothing
	return nil
}

// loadConfigFile load configuration file and replace ENV data in configuration file
// priority: flag > env > configuration. If ENV not existed, take flag default value
func loadConfigFile(app *App, cmd *cobra.Command) error {
	viper.SetConfigFile(app.cfOption.path) // need by viper.ReadConfig()

	buf, err := ioutil.ReadFile(app.cfOption.path)
	if err != nil {
		return err
	}

	buf, err = replaceENV(buf)
	if err != nil {
		return err
	}
	//dumpConfigBuf(buf, "after replace by ENV")

	if err := viper.ReadConfig(bytes.NewReader(buf)); err != nil {
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

/* support for ENV variables */
// format: MSENV(ALL_SECTION_KEY); all capital
const MARK_LEFT = "MSENV("
const MARK_RIGHT = ")"

func replaceENV(buf []byte) ([]byte, error) {
	data := make([]byte, 0, len(buf))

	scanner := bufio.NewScanner(bytes.NewReader(buf))
	for scanner.Scan() {
		line := append(scanner.Bytes(), '\n')

		left := bytes.Index(line, []byte(MARK_LEFT))
		if left == -1 {
			data = append(data, line...)
			continue
		}

		right := bytes.Index(line[left:], []byte(MARK_RIGHT))
		if right == -1 {
			data = append(data, line...)
			continue
		}
		right += left // right should be the index of hold line

		/* replace with ENV */
		placeholder := line[left : right+1]
		envStr := string(line[left+len(MARK_LEFT) : right])
		envData := os.Getenv(envStr)
		line = bytes.Replace(line, placeholder, []byte(envData), 1)

		data = append(data, line...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func dumpConfigBuf(buf []byte, desc string) {
	fmt.Printf("\n\n>> ========== %s ========== <<\n", desc)
	fmt.Println(string(buf))
	fmt.Printf(">> ================================ <<\n\n")
}
