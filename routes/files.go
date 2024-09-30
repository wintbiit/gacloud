package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/gacloud/model"
	"github.com/wintbiit/gacloud/server"
	"strconv"
)

func init() {
	addHookAuth("/api/v1/files", func(party iris.Party) {
		party.Get("/list", List)
		party.Post("/upload", Upload)
		party.Get("/get", Get)
	})

	addHook("/dl", func(party iris.Party) {
		party.Get("/:sum", DownloadFile)
	})
}

// TODO: Pagnation
func List(ctx iris.Context) {
	path := ctx.URLParam("path")
	user := jwt.Get(ctx).(*model.UserClaims).ToUser()
	s := server.GetServer()

	list, clean, err := s.ListFiles(ctx, user, path)
	defer clean()
	if err != nil {
		log.Err(err).Msg("Failed to list files")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(list)
}

func Upload(ctx iris.Context) {
	path := ctx.URLParam("path")
	user := jwt.Get(ctx).(*model.UserClaims).ToUser()
	s := server.GetServer()

	// 1. Upload API does not overwrite existing files, so if the file exists, throw
	exists, err := s.FileExists(ctx, user, path)
	if err != nil {
		log.Err(err).Msg("Failed to check file exists")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	if exists {
		ctx.StopWithStatus(iris.StatusConflict)
		return
	}

	// 2. Get file size
	lengthStr := ctx.GetHeader("Content-Length")
	if lengthStr == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	length, err := strconv.ParseUint(lengthStr, 10, 64)
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// 3. Get user default provider
	provider := s.GetDefaultProviderID(ctx, user)

	f, err := model.NewFile(user, path, length, "", provider)

	// 4. Get file content, calculate checksum
	reader := ctx.Request().Body

	checksum, err := s.WriteFile(ctx, user, f, reader)
	if err != nil {
		log.Err(err).Msg("Failed to write file")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	f.Sum = checksum

	// 5. Put file metadata to elasticsearch
	err = s.PutFile(ctx, user, f)
	if err != nil {
		log.Err(err).Msg("Failed to put file metadata")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	ctx.StopWithStatus(iris.StatusOK)
}

type getFileResponse struct {
	File *model.File `json:"file"`
	Sign string      `json:"sign"`
}

func Get(ctx iris.Context) {
	path := ctx.URLParam("path")
	sum := ctx.URLParam("sum")
	user := jwt.Get(ctx).(*model.UserClaims).ToUser()
	s := server.GetServer()

	var file *model.File
	var err error
	if path != "" {
		file, _, err = s.GetFileBySum(ctx, user, path)
		if err != nil {
			log.Err(err).Msg("Failed to get file")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}
	} else if sum != "" {
		file, _, err = s.GetFileBySum(ctx, user, sum)
		if err != nil {
			log.Err(err).Msg("Failed to get file")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}
	} else {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	if file == nil {
		ctx.StopWithStatus(iris.StatusNotFound)
		return
	}

	sign, err := s.GenerateDownloadSign(ctx, user, file)
	if err != nil {
		log.Err(err).Msg("Failed to generate download sign")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(getFileResponse{File: file, Sign: sign})
}

func DownloadFile(ctx iris.Context) {
	sum := ctx.Params().Get("sum")
	sign := ctx.URLParam("sign")
	s := server.GetServer()

	if sign == "" {
		ctx.StopWithStatus(iris.StatusForbidden)
		return
	}

	if !s.VerifyDownloadSign(ctx, sum, sign) {
		ctx.StopWithStatus(iris.StatusForbidden)
		return
	}

	file, clean, err := s.GetFileBySumNoCheck(ctx, sum)
	if err != nil {
		log.Err(err).Msg("Failed to get file")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	defer clean()

	if file == nil {
		ctx.StopWithStatus(iris.StatusNotFound)
		return
	}

	provider, err := s.GetProvider(file.ProviderId)
	if err != nil {
		log.Err(err).Msg("Failed to get provider")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	reader, ok, err := provider.Get(ctx, file.Sum)
	if err != nil {
		log.Err(err).Msg("Failed to get file")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	if !ok {
		ctx.StopWithStatus(iris.StatusNotFound)
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename="+file.Name())
	ctx.Header("Content-Type", file.Mime)
	ctx.ServeContent(reader, file.Name(), file.UpdatedAt)
}
