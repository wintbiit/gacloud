package model

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"gorm.io/gorm"
)

const (
	FileOwnerTypeUser   = 0
	FileOwnerTypeGroup  = 1
	FileOwnerTypeShared = 2
)

type File struct {
	Path       string `json:"path" `
	Size       int64  `json:"size"`
	Mime       string `json:"mime" `
	OwnerType  int8   `json:"owner_type" `
	OwnerId    uint   `json:"owner_id"`
	Sum        string `json:"sum"`
	ProviderId uint   `json:"provider_id"`
	Fp         string `json:"fp,omitempty"`
}

var FileTypeMapping = &types.TypeMapping{
	Properties: map[string]types.Property{
		"sum":         types.NewTextProperty(),
		"path":        types.NewKeywordProperty(),
		"size":        types.NewIntegerNumberProperty(),
		"mime":        types.NewTextProperty(),
		"owner_type":  types.NewIntegerNumberProperty(),
		"owner_id":    types.NewIntegerNumberProperty(),
		"provider_id": types.NewIntegerNumberProperty(),
	},
}

type FileProvider struct {
	gorm.Model
	Name       string `gorm:"unique"`
	Type       string `gorm:"not null"`
	Credential string `gorm:"not null"`
}
