package command

import (
	"MarineSnow/framework/cobra"
)

func init() {
	RootCmd.AddCommand(appCmd)
}

var RootCmd = &cobra.Command{
	Use:               "MarineSnow",
	Short:             "MarineSnow web framework CLI",
	Long:              `A web framework command line tool`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		// do nothing in "MarineSnow" command
		cmd.InitDefaultHelpFlag()
		return cmd.Help()
	},
}
