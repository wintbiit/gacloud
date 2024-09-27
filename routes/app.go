package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/wintbiit/gacloud/server"
)

func init() {
	addHook("/api/v1", func(app iris.Party) {
		app.Get("/appinfo", AppInfo)
	})
}

func AppInfo(ctx iris.Context) {
	ctx.JSON(server.GetServer().Info)
}
