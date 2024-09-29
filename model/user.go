package model

import (
	"path"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string `gorm:"unique,not null,index"`
	DisplayName string `gorm:"not null"`
	Email       string `gorm:"unique,not null,index"`
	Password    string `gorm:"not null"`
}

func (u *User) HomeDir() string {
	return path.Join(UserScopeDir, u.Name)
}
