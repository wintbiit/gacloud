package model

const FileOwnerTypeUser = "user"
const FileOwnerTypeGroup = "group"
const FileOwnerTypeShared = "shared"

type File struct {
	Path       string `xorm:"varchar(255) notnull"`
	Size       int64  `xorm:"notnull"`
	Sum        string `xorm:"varchar(64) notnull"`
	Mime       string `xorm:"varchar(255) notnull"`
	ProviderId int64  `xorm:"notnull"`
	TimeModel
	OwnerModel
}

type FileProvider struct {
	Name       string `xorm:"varchar(255) notnull"`
	Type       string `xorm:"varchar(255) notnull"`
	Credential string `xorm:"json notnull"`
	TimeModel
}
