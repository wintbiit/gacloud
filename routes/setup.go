package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/wintbiit/gacloud/internal"
)

func init() {
	addHook("/setup", false, func(app iris.Party) {
		app.Get("/config", GetSetupConfig)
		app.Post("/database", SetDataBaseOptions)
	})
}

func GetSetupConfig(ctx iris.Context) {

}

type DataBaseOption struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}

func SetDataBaseOptions(ctx iris.Context) {
	var dbOption DataBaseOption
	if err := ctx.ReadJSON(&dbOption); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	if err := internal.TryDbOptions(dbOption.Type, dbOption.Params); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	ctx.StatusCode(iris.StatusOK)
}
