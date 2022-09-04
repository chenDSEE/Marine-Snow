package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"MarineSnow/framework/cobra"
	"MarineSnow/framework/system"

	"github.com/robfig/cron/v3"
	"github.com/sevlyar/go-daemon"
)

var isDaemon bool

func init() {
	// TODO: set the interval for cron-job by command
	// TODO: MarineSnow cron start [cron-job-name]
	// Should not hard code all cron-job in cronStartCmd
	//
	cronStartCmd.Flags().BoolVarP(&isDaemon, "daemon", "d", false, "start cron job as a daemon")
	cronCmd.AddCommand(cronStartCmd)
	cronCmd.AddCommand(cronStopCmd)
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

// TODO: specified by app, but not hard code here
const workDir = "./"
const pidFileName = "storage/runtime/cron.pid"
const logFileName = "storage/log/cron.log"

var cronStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start a cron job",
	Long:  "start a cron job in MarineSnow",
	RunE: func(cmd *cobra.Command, args []string) error {
		// register CPU resource
		runtime.GOMAXPROCS(1)

		/* init Cron and cron job */
		c := cron.New(cron.WithSeconds())
		_, err := c.AddFunc("* * * * * *", func() {
			MonitorSystemStatus()
		})

		if err != nil {
			fmt.Printf("cron job fail to start, error:%s\n", err.Error())
			return err
		}

		pid := os.Getpid()
		/* start cron job */
		if isDaemon { // cron job running as a daemon
			/* setup daemon variable */
			dCtx := &daemon.Context{
				PidFileName: pidFileName,
				PidFilePerm: 0644,
				LogFileName: logFileName,
				LogFilePerm: 0640,
				WorkDir:     workDir,
				Umask:       027,
				Args:        os.Args, // make cobra run correctly
			}

			/* create and run daemon */
			child, err := dCtx.Reborn()
			if err != nil {
				fmt.Println("Create child process failed: ", err)
				return err
			}

			if child != nil {
				// parent process
				printCronStartInfo(isDaemon, child.Pid)
				return nil
			}

			// child process
			defer dCtx.Release()
			pid = os.Getpid() // update child pid

		} else { // cron job running as a foreground
			err = ioutil.WriteFile(workDir+pidFileName, []byte(strconv.Itoa(pid)), 0664)
			if err != nil {
				fmt.Printf("failed to update pid file, err: %s\n", err.Error())
			}
		}

		printCronStartInfo(isDaemon, pid)
		c.Run()

		return nil
	},
}

func printCronStartInfo(isDaemon bool, pid int) {
	line := "==================="
	title := "cron job start"
	replace := make([]uint8, 256, 256)
	for i := 0; i < len(replace) && i < len(title); i++ {
		replace[i] = '='
	}

	fmt.Printf("+%s %s %s+\n", line, title, line)
	fmt.Printf("|start time [%s]\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("|PID:%d, stored in [%s]\n", pid, workDir+pidFileName)

	if isDaemon {
		fmt.Printf("|Log stored in [%s]\n", workDir+logFileName)
	}

	fmt.Printf("+%s=%s=%s+\n", line, string(replace), line)
}

func MonitorSystemStatus() {
	fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), system.Info())
}

var cronStopCmd = &cobra.Command{
	Use:              "stop",
	Short:            "stop a cron job",
	Long:             "stop a cron job in MarineSnow",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := ioutil.ReadFile(workDir + pidFileName)
		if err != nil {
			fmt.Printf("fail to get cron jon pid from %s, err: %s\n",
				workDir+pidFileName, err.Error())

			return err
		}

		if len(data) <= 0 {
			return err
		}

		pid, err := strconv.Atoi(string(data))
		if err != nil {
			fmt.Printf("fail to convert string to number\n")
			return err
		}

		// kill cron jon
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			return err
		}

		// reset pid-file
		if err := ioutil.WriteFile(workDir+pidFileName, []byte{}, 0644); err != nil {
			return err
		}

		fmt.Printf("[%s] stop cron job [PID:%d] success\n",
			time.Now().Format("2006-01-02 15:04:05"), pid)

		return nil
	},
}
