package cmd

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/wintbiit/gacloud/routes"
)

func Daemon() {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	routes.RegisterRoutes(app.Party("/api/v1"))
}
