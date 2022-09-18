package command

import (
	"MarineSnow/app/appDemo"
	"MarineSnow/framework/cobra"
	"MarineSnow/framework/config/env"
	"fmt"
)

var AppName string
var ipAddr string
var port string

func init() {
	appCmd.PersistentFlags().StringVarP(&AppName, "name", "n", "", "the app name")
	appCmd.PersistentFlags().StringVarP(&ipAddr, "ip", "i", "127.0.0.1", "ip addr listen on")
	appCmd.PersistentFlags().StringVarP(&port, "port", "p", "80", "port listen on")
	appCmd.AddCommand(appStartCmd)
	//appCmd.AddCommand(appstopCmd) TODO:
}

var appCmd = &cobra.Command{
	Use:              "app",
	Short:            "app command",
	Long:             "MarineSnow app command",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// do nothing in "MarineSnow app" command
		return cmd.Help()
	},
}

var appStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start a app",
	Long:  "start a app register in MarineSnow",
	RunE: func(cmd *cobra.Command, args []string) error {
		/* init configuration */
		env.EnvInit("./config/")

		/* try to start app */
		fmt.Printf("try to run app and start it, addr: %s, port: %s\n", ipAddr, port)
		appDemo.StartAppDemo(ipAddr + ":" + port)
		return nil
	},
}
