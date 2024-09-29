package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/server"
)

func init() {
	addHookAuth("/api/v1/files", func(party iris.Party) {
		party.Get("/:path", List)
	})
}

// TODO: Pagnation
func List(ctx iris.Context) {
	path := ctx.Params().Get("path")
	user := jwt.Get(ctx).(*model.UserClaims).ToUser()
	s := server.GetServer()

	list, clean, err := s.ListFiles(ctx, user, path)
	defer clean()
	if err != nil {
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(list)
}
