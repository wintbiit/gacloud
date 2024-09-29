package model

import (
	"path"

	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Name        string `gorm:"unique,not null"`
	DisplayName string `gorm:"not null"`
}

type UserGroup struct {
	gorm.Model
	User    User `gorm:"foreignKey:UserID"`
	UserID  uint
	Group   Group `gorm:"foreignKey:GroupID"`
	GroupID uint
}

func (g *Group) HomeDir() string {
	return path.Join(GroupScopeDir, g.Name)
}

type UserFileProvider struct {
	gorm.Model
	User           User `gorm:"foreignKey:UserID"`
	UserID         uint
	FileProvider   FileProvider `gorm:"foreignKey:FileProviderID"`
	FileProviderID uint
}

type GroupFileProvider struct {
	gorm.Model
	Group          Group `gorm:"foreignKey:GroupID"`
	GroupID        uint
	FileProvider   FileProvider `gorm:"foreignKey:FileProviderID"`
	FileProviderID uint
}
