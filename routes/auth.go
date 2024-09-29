package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/gacloud/server"
)

func init() {
	addHook("/api/v1", func(party iris.Party) {
		party.Post("/login", Login)
	})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(ctx iris.Context) {
	var req LoginRequest
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	s := server.GetServer()

	user, err := s.UserLogin(ctx.Request().Context(), req.Username, req.Password)
	if err != nil {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		return
	}

	token, err := s.GenerateUserToken(user)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate token")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(LoginResponse{Token: token})
}
