package model

import "gorm.io/gorm"

const (
	FileOwnerTypeUser   = 0
	FileOwnerTypeGroup  = 1
	FileOwnerTypeShared = 2
)

type File struct {
	Path       string `json:"path" gorm:"unique,not null,index"`
	Size       int64  `json:"size" gorm:"not null"`
	Mime       string `json:"mime" gorm:"not null"`
	OwnerType  int8   `json:"owner_type" gorm:"not null"`
	OwnerId    uint   `json:"owner_id" gorm:"not null"`
	Sum        string `gorm:"unique,not null,index" json:"sum"`
	ProviderId uint   `gorm:"not null" json:"provider_id"`
	Fd         string `gorm:"-" json:"fd"`
}

type FileProvider struct {
	gorm.Model
	Name       string `gorm:"unique"`
	Type       string `gorm:"not null"`
	Credential string `gorm:"not null"`
}
