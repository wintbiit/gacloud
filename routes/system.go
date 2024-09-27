package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/wintbiit/gacloud/server"
	"github.com/wintbiit/gacloud/utils"
)

func init() {
	addHookFront("/api/v1", func(app iris.Party) {
		app.Get("/health", Health)
		app.Get("/serverinfo", ServerInfo)
	})
}

func Health(ctx iris.Context) {
	if server.GetServer() == nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	ctx.StatusCode(iris.StatusOK)
}

func ServerInfo(ctx iris.Context) {
	ctx.JSON(utils.ServerInfo)
}
