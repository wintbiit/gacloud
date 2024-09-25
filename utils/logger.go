package utils

import (
	"github.com/kataras/iris/v12"
	"os"
	"strconv"
	"time"
	xormlog "xorm.io/xorm/log"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var DEBUG = os.Getenv("GACLOUD_DEBUG") == "true"

func init() {
	lumberjackLogger := lumberjack.Logger{
		Filename:   "logs/auther.log",
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
		Compress:   true,
	}

	log.Logger = log.Output(zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{Out: os.Stderr},
		zerolog.SyncWriter(&lumberjackLogger),
	))

	if DEBUG {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Msg("Logger initialized")
}

func NewIrisLogger() iris.Handler {
	logger := log.Logger.With().Str("framework", "iris").Logger()

	return func(ctx iris.Context) {
		var status, ip, method, path string
		var latency time.Duration
		var startTime, endTime time.Time
		startTime = time.Now()

		// Before Next.
		ip = ctx.RemoteAddr()
		method = ctx.Method()

		ctx.Next()

		endTime = time.Now()
		latency = endTime.Sub(startTime)
		status = strconv.Itoa(ctx.GetStatusCode())
		path = ctx.Path()

		logger.Info().
			Str("ip", ip).
			Str("method", method).
			Str("path", path).
			Str("status", status).
			Dur("latency", latency).
			Msg("request finished")
	}
}

type xormLogger struct {
	*zerolog.Logger
	level   xormlog.LogLevel
	showSQL bool
}

func (x xormLogger) Debug(v ...interface{}) {
	x.Logger.Debug().Msg("xorm debug")
}

func (x xormLogger) Debugf(format string, v ...interface{}) {
	x.Logger.Debug().Msgf(format, v...)
}

func (x xormLogger) Error(v ...interface{}) {
	x.Logger.Error().Msg("xorm error")
}

func (x xormLogger) Errorf(format string, v ...interface{}) {
	x.Logger.Error().Msgf(format, v...)
}

func (x xormLogger) Info(v ...interface{}) {
	x.Logger.Info().Msg("xorm info")
}

func (x xormLogger) Infof(format string, v ...interface{}) {
	x.Logger.Info().Msgf(format, v...)
}

func (x xormLogger) Warn(v ...interface{}) {
	x.Logger.Warn().Msg("xorm warn")
}

func (x xormLogger) Warnf(format string, v ...interface{}) {
	x.Logger.Warn().Msgf(format, v...)
}

func (x xormLogger) Level() xormlog.LogLevel {
	return x.level
}

func (x xormLogger) SetLevel(l xormlog.LogLevel) {
	x.level = l
}

func (x xormLogger) ShowSQL(show ...bool) {
	if len(show) > 0 {
		x.showSQL = show[0]
	}
}

func (x xormLogger) IsShowSQL() bool {
	return x.showSQL
}

func NewXormLogger() xormlog.Logger {
	logger := log.Logger.With().Str("framework", "xorm").Logger()
	level := xormlog.LOG_INFO
	if DEBUG {
		level = xormlog.LOG_DEBUG
	}

	showSql := DEBUG

	return &xormLogger{
		Logger:  &logger,
		level:   level,
		showSQL: showSql,
	}
}
