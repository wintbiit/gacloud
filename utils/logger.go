package utils

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	gormlog "gorm.io/gorm/logger"
)

var (
	DEBUG  = os.Getenv("GACLOUD_DEBUG") == "true"
	logDir = GetEnv("GACLOUD_LOG_DIR", "./logs")
)

func init() {
	lumberjackLogger := lumberjack.Logger{
		Filename:   logDir + "/gacloud.log",
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

func NewLogger(framework string) zerolog.Logger {
	return log.With().Str("framework", framework).Logger()
}

func NewIrisLogger() iris.Handler {
	logger := NewLogger("iris")

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
			Stringer("latency", latency).
			Str("user_id", ctx.Values().GetString("user_id")).
			Send()
	}
}

type gormLogger struct {
	*zerolog.Logger
	level gormlog.LogLevel
}

func (g *gormLogger) LogMode(level gormlog.LogLevel) gormlog.Interface {
	g.level = level
	return g
}

func (g *gormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	g.Logger.Info().Ctx(ctx).Msgf(s, i...)
}

func (g *gormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	g.Logger.Warn().Ctx(ctx).Msgf(s, i...)
}

func (g *gormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	g.Logger.Error().Ctx(ctx).Msgf(s, i...)
}

func (g *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil:
		sql, rows := fc()
		g.Logger.Error().
			Ctx(ctx).
			Str("sql", sql).
			Int64("rows", rows).
			Err(err).
			Dur("elapsed", elapsed).
			Send()
	default:
		sql, rows := fc()
		g.Logger.Debug().
			Ctx(ctx).
			Str("sql", sql).
			Int64("rows", rows).
			Dur("elapsed", elapsed).
			Send()
	}
}

func NewGormLogger() gormlog.Interface {
	logger := NewLogger("gorm")
	level := gormlog.Warn
	if DEBUG {
		level = gormlog.Info
	}
	return &gormLogger{Logger: &logger, level: level}
}
