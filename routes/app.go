package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/wintbiit/gacloud/config"
)

func init() {
	addHook("/", true, func(app iris.Party) {
		app.Get("/info", GetAppInfo)
	})
}

type AppInfo struct {
	SiteName    string `json:"site_name"`
	ExternalUrl string `json:"external_url"`
	SiteIcon    string `json:"site_icon"`
}

var appInfo *AppInfo

func GetAppInfo(ctx iris.Context) {
	if appInfo == nil {
		appInfo = &AppInfo{
			SiteName:    config.GetWithDefault("app.name", "GaCloud"),
			ExternalUrl: config.GetWithDefault("app.external_url", "http://localhost:8080"),
			SiteIcon:    config.GetWithDefault("app.icon", ""),
		}
	}

	ctx.JSON(appInfo)
}
