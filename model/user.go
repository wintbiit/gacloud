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

func (u *User) ToClaims() *UserClaims {
	return &UserClaims{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

type UserClaims struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *UserClaims) ToUser() *User {
	return &User{
		Model: gorm.Model{ID: u.ID},
		Name:  u.Name,
		Email: u.Email,
	}
}
