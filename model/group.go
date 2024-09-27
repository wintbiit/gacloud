package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name        string `gorm:"unique,not null"`
	DisplayName string `gorm:"not null"`
}

type UserGroup struct {
	gorm.Model
	User  *User  `gorm:"foreignKey:UserID"`
	Group *Group `gorm:"foreignKey:GroupID"`
}
