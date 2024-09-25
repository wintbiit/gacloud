package model

type Group struct {
	Name        string `xorm:"varchar(25) notnull unique 'name'"`
	DisplayName string `xorm:"varchar(25) notnull 'display_name'"`
	TimeModel
}

type UserGroup struct {
	UserId  int64 `xorm:"varchar(36) notnull"`
	GroupId int64 `xorm:"varchar(36) notnull"`
	TimeModel
}
