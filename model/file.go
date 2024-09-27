package model

import "gorm.io/gorm"

const (
	FileOwnerTypeUser   = 0
	FileOwnerTypeGroup  = 1
	FileOwnerTypeShared = 2
)

type ListFile struct {
	Path      string `json:"path" gorm:"unique,not null,index"`
	Size      int64  `json:"size" gorm:"not null"`
	Mime      string `json:"mime" gorm:"not null"`
	OwnerType int8   `json:"owner_type" gorm:"not null"`
	OwnerId   uint   `json:"owner_id" gorm:"not null"`
}

type File struct {
	ListFile
	Sum        string `gorm:"unique,not null,index" json:"sum"`
	ProviderId int64  `gorm:"not null" json:"provider_id"`
}

type FileProvider struct {
	gorm.Model
	Name       string `gorm:"unique"`
	Type       string `gorm:"not null"`
	Credential string `gorm:"not null"`
}
