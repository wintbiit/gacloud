package internal

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/gacloud/model"
)

const FilePermRD = 1
const FilePermRW = 2
const FilePermRMN = 3

func (s *GaCloudServer) AuthorizeFileAction(ctx context.Context, user *model.User, file *model.File, action int) bool {
	switch file.OwnerType {
	case model.FileOwnerTypeUser:
		return file.OwnerId == user.ID
	case model.FileOwnerTypeGroup:
		return s.IsUserInGroup(ctx, user, file.OwnerId)
	default:
		log.Error().Int64("fileId", file.ID).Str("ownerType", file.OwnerType).Msg("Unknown file owner type")
		return false
	}
}
