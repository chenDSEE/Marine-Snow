package main

import (
	SnowCommand "MarineSnow/framework/command"
	"fmt"
)

func main() {
	fmt.Printf("Welcome to MarineSnow, framework start now.\n")

	if err := SnowCommand.RootCmd.Execute(); err != nil {
		fmt.Printf("Run command error, %s\n", err.Error())
		return
	}

	fmt.Println("Bye~, exit Marine-Snow now.")
}
