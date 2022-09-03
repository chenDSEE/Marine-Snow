package command

import (
	"MarineSnow/framework/cobra"
	"MarineSnow/framework/system"
	"fmt"
	"github.com/robfig/cron/v3"
	"runtime"
	"time"
)

func init() {
	// TODO: set the interval for cron-job by command
	// TODO: MarineSnow cron start [cron-job-name]
	// Should not hard code all cron-job in cronStartCmd
	//
	cronCmd.AddCommand(cronStartCmd)
	//cronCmd.AddCommand(cronStopCmd)
	//cronCmd.AddCommand(cronStatusCmd)
}

// TODO: wrap cron command as a option, register cron job when init MarineSnow framework
var cronCmd = &cobra.Command{
	Use:              "cron",
	Short:            "cron command",
	Long:             "MarineSnow cron job command",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// do nothing in "MarineSnow cron" command
		return cmd.Help()
	},
}

var cronStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start a cron job",
	Long:  "start a cron job in MarineSnow",
	RunE: func(cmd *cobra.Command, args []string) error {
		// register CPU resource
		runtime.GOMAXPROCS(1)

		/* init Cron and cron job */
		c := cron.New(cron.WithSeconds())
		c.AddFunc("* * * * * *", func() {
			MonitorSystemStatus()
		})

		/* start cron job */
		c.Run()
		return nil
	},
}

func MonitorSystemStatus() {
	fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), system.Info())
}
