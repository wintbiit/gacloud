package model

type User struct {
	Name        string `xorm:"varchar(25) notnull unique 'name'"`
	DisplayName string `xorm:"varchar(25) notnull 'display_name'"`
	Email       string `xorm:"varchar(50) notnull unique 'email'"`
	Password    string `xorm:"varchar(50) notnull 'password'"`
	TimeModel
}
