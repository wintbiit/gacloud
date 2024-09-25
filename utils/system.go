package utils

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

var Version = "v0.0.1"

var shutdownHooks []func()

func init() {
	shutdownHooks = make([]func(), 0)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Warn().Str("signal", sig.String()).Msg("GaCloud is terminating")
		for _, hook := range shutdownHooks {
			hook()
		}
	}()
}

func GetVersion() {
}

func AddShutdownHook(hook func()) {
	shutdownHooks = append(shutdownHooks, hook)
}
