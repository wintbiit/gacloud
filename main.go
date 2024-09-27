package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/wintbiit/gacloud/cmd"
	"github.com/wintbiit/gacloud/utils"
)

const Banner = `
GaCloud - A cloud storage service aims to provide group access control and openid connect support.

Version: %s
Go Version: %s
Build Time: %s
Commit Hash: %s
`

const (
	DebugMsg = `Debug mode enabled`
	HelpMsg  = `Usage`
)

func main() {
	args := os.Args

	fmt.Printf(Banner, utils.ServerInfo.Version, utils.ServerInfo.GoVersion, utils.ServerInfo.BuildTime, utils.ServerInfo.BuildRevision)
	if utils.DEBUG {
		fmt.Println(DebugMsg)
	}

	if len(args) < 2 {
		fmt.Println()
		os.Exit(1)
	}

	switch args[1] {
	case "version":
		log.Info().Interface("version", utils.ServerInfo).Msg("GaCloud Version")
	case "help":
		fmt.Println(HelpMsg)
	case "daemon":
		cmd.Daemon()
	default:
		fmt.Println(HelpMsg)
		os.Exit(1)
	}
}
