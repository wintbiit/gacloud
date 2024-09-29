package model

import (
	"path"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string `gorm:"unique,not null"`
	DisplayName string `gorm:"not null"`
	Email       string `gorm:"unique,not null"`
	Password    string `gorm:"not null"`
}

func (u *User) HomeDir() string {
	return path.Join(UserScopeDir, u.Name)
}
