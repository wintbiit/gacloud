package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/wintbiit/gacloud/internal"
)

func init() {
	addHook("/", false, func(app iris.Party) {
		app.Get("/health", Healthz)
	})
}

func Healthz(ctx iris.Context) {
	if internal.GetGaCloudServer() == nil || !internal.GetGaCloudServer().HealthCheck() {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString("Database connection failed")
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.WriteString("OK")
}
