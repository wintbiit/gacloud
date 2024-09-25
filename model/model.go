package model

import "time"

type TimeModel struct {
	ID        int64      `xorm:"pk autoincr"`
	CreatedAt time.Time  `xorm:"created"`
	UpdatedAt time.Time  `xorm:"updated"`
	DeleteAt  *time.Time `xorm:"deleted"`
}

type OwnerModel struct {
	OwnerId   int64  `xorm:"varchar(36) notnull"`
	OwnerType string `xorm:"varchar(36) notnull"`
}
