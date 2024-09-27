package cmd

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/gacloud/routes"
	"github.com/wintbiit/gacloud/server"
	"github.com/wintbiit/gacloud/utils"
)

func Daemon() {
	app := iris.New()
	app.Use(recover.New())
	app.Use(utils.NewIrisLogger())

	app.HandleDir("/public", "./web/dist")

	routes.RegisterRoutes(app)

	s := server.GetServer()
	if s == nil {
		log.Warn().Msg("server is not setup")
	} else {
		log.Info().Interface("server info", s.Info).Msg("local server")
	}

	log.Info().Str("addr", utils.ServerInfo.Addr).Msg("starting server")
	if err := app.Listen(utils.ServerInfo.Addr,
		iris.WithCharset("UTF-8"),
		iris.WithOtherValue("ServerName", "GaCloud"),
		iris.WithEmptyFormError,
		iris.WithPostMaxMemory(32<<20),
		iris.WithoutBanner,
		iris.WithEasyJSON,
		iris.WithoutStartupLog,
	); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
