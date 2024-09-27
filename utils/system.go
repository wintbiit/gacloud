package utils

import (
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"syscall"

	"github.com/rs/zerolog/log"
)

type ServerInformation struct {
	Version       string `json:"version"`
	BuildRevision string `json:"build_revision"`
	BuildTime     string `json:"build_time"`
	GoVersion     string `json:"go_version"`
	DataDir       string `json:"-"`
	LogDir        string `json:"-"`
	Addr          string `json:"addr"`
}

var (
	shutdownHooks []func()
	version       = "unknown"
	ServerInfo    *ServerInformation
)

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

	dataDir := GetEnv("GACLOUD_DATA_DIR", "./")
	logDir := GetEnv("GACLOUD_LOG_DIR", "./logs")

	// to abs path
	dataDir, _ = filepath.Abs(dataDir)
	logDir, _ = filepath.Abs(logDir)

	os.MkdirAll(dataDir, 0o755)
	os.MkdirAll(logDir, 0o755)

	ServerInfo = &ServerInformation{
		Version:   version,
		GoVersion: runtime.Version(),
		DataDir:   dataDir,
		LogDir:    logDir,
		Addr:      GetEnv("GACLOUD_ADDR", ":8080"),
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, setting := range buildInfo.Settings {
			if version != "" {
				break
			}

			if setting.Key == "vcs.revision" {
				ServerInfo.BuildRevision = setting.Value[0:7]
			}

			if setting.Key == "vcs.time" {
				ServerInfo.BuildTime = setting.Value
			}
		}
	}
}

func AddShutdownHook(hook func()) {
	shutdownHooks = append(shutdownHooks, hook)
}
