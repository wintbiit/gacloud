package routes

import (
	"github.com/kataras/iris/v12"
)

func init() {
	addHookFront("/api/v1/setup", func(app iris.Party) {
		app.Get("/", GetSetupStatus)
		app.Post("/database", SetDataBaseOptions)
		app.Post("/database/test", TestDataBase)
	})
}

type SetupStatusResponse struct {
	CurrentStep int `json:"currentStep"`
}

func GetSetupStatus(ctx iris.Context) {
	ctx.JSON(SetupStatusResponse{
		CurrentStep: 1,
	})
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

	ctx.StatusCode(iris.StatusOK)
}

type TestDataBaseResponse struct {
	Success bool `json:"success"`
}

func TestDataBase(ctx iris.Context) {
	var dbOption DataBaseOption
	if err := ctx.ReadJSON(&dbOption); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	ctx.JSON(TestDataBaseResponse{
		Success: true,
	})
}
