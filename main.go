package main

import (
	"fmt"
	"github.com/wintbiit/gacloud/cmd"
	"github.com/wintbiit/gacloud/utils"
	"os"
)

const Banner = `
GaCloud - A cloud storage service aims to provide group access control and openid connect support.

Version: %s
`
const DebugMsg = `Debug mode enabled`
const HelpMsg = `Usage`

func main() {
	args := os.Args

	fmt.Printf(Banner, utils.Version)
	if utils.DEBUG {
		fmt.Println(DebugMsg)
	}

	if len(args) < 2 {
		fmt.Println(HelpMsg)
		os.Exit(1)
	}

	switch args[1] {
	case "version":
		utils.GetVersion()
	case "help":
		fmt.Println(HelpMsg)
	case "daemon":
		cmd.Daemon()
	default:
		fmt.Println(HelpMsg)
		os.Exit(1)
	}
}
