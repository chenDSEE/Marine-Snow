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
	cronCmd.AddCommand(cronRestartCmd)
	cronCmd.AddCommand(cronStatusCmd)
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
		return startCronJob(cmd, args, os.Args)
	},
}

func startCronJob(cmd *cobra.Command, args []string, osArgs []string) error {
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
			Args:        osArgs, // make cobra run correctly
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

// only support for cron job daemon
var cronRestartCmd = &cobra.Command{
	Use:              "restart",
	Short:            "restart a cron job daemon",
	Long:             "restart a cron job daemon in MarineSnow",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// get old cron job pid
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

		// try to kill old cron job in 10 times with 1s interval
		if isProcessExisted(pid) {
			_ = syscall.Kill(pid, syscall.SIGTERM)

			times := 0
			for ; times < 10; times++ {
				if isProcessExisted(pid) == false {
					break
				}

				if times != 0 && times/2 == 0 {
					_ = syscall.Kill(pid, syscall.SIGTERM)
					fmt.Printf("fail to kill PID:%d, try again after 2s.\n", pid)
				}

				time.Sleep(1 * time.Second)
			}

			if times == 10 {
				fmt.Printf("fail to kill PID:%d\n", pid)
				return nil
			}

			fmt.Printf("[%s] stop cron job [PID:%d] success\n",
				time.Now().Format("2006-01-02 15:04:05"), pid)
		}

		// reset pid-file
		if err := ioutil.WriteFile(workDir+pidFileName, []byte{}, 0644); err != nil {
			return err
		}

		// start again as a daemon
		newArgs := make([]string, 0, len(os.Args)+1)
		for _, arg := range os.Args {
			if arg == "restart" {
				newArgs = append(newArgs, "start")
				continue
			}

			newArgs = append(newArgs, arg)
		}

		newArgs = append(newArgs, "-d")
		isDaemon = true
		return startCronJob(cmd, args, newArgs)
	},
}

var cronStatusCmd = &cobra.Command{
	Use:              "status",
	Short:            "show the status of cron job",
	Long:             "show the status of cron job in MarineSnow",
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

		// check process status via pid
		isExisted := isProcessExisted(pid)
		status := "Running"
		if isExisted == false {
			status = "Stop"
		}

		printCronStatus(pid, status)

		return nil
	},
}

// return true when process existed
func isProcessExisted(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// check process alive via signal
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}

	return true
}

func printCronStatus(pid int, status string) {
	line := "==================="
	title := "cron status"
	replace := make([]uint8, 256, 256)
	for i := 0; i < len(replace) && i < len(title); i++ {
		replace[i] = '='
	}

	fmt.Printf("+%s %s %s+\n", line, title, line)

	fmt.Printf("|PID[%d] %s\n", pid, status)

	fmt.Printf("+%s=%s=%s+\n", line, string(replace), line)
}
