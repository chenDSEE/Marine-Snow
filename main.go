package main

import (
	SnowCommand "MarineSnow/framework/command"
	"fmt"
)

func main() {
	if err := SnowCommand.RootCmd.Execute(); err != nil {
		fmt.Printf("Run command error, %s\n", err.Error())
		return
	}

	fmt.Printf("appName[%s]\n", SnowCommand.AppName)
}
