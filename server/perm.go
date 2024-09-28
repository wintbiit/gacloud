package server

import (
	"context"

	"github.com/wintbiit/gacloud/model"
)

const (
	FilePermRD  = 1
	FilePermRW  = 2
	FilePermRMN = 3
)

func (s *GaCloudServer) AuthorizeFileAction(ctx context.Context, user *model.User, file *model.File, action int) bool {
	switch file.OwnerType {
	case model.FileOwnerTypeUser:
		return file.OwnerId == user.ID
	case model.FileOwnerTypeGroup:
		return s.IsUserInGroup(ctx, user, file.OwnerId)
	default:
		s.logger.Error().Str("sum", file.Sum).
			Int8("ownerType", file.OwnerType).
			Uint("ownerId", file.OwnerId).
			Msg("Unknown file owner type")
		return false
	}
}
