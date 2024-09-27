package internal

import (
	"github.com/wintbiit/gacloud/utils"
)

func init() {
	utils.AddShutdownHook(func() {
	})
}
