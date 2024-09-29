package model

import (
	"context"

	"gorm.io/gorm"
)

func MigrateModels(ctx context.Context, db *gorm.DB) error {
	err := db.AutoMigrate(&File{})
	if err != nil {
		return err
	}

	err = db.WithContext(ctx).AutoMigrate(&FileProvider{})
	if err != nil {
		return err
	}

	err = db.WithContext(ctx).AutoMigrate(&Group{})
	if err != nil {
		return err
	}

	err = db.WithContext(ctx).AutoMigrate(&UserGroup{})
	if err != nil {
		return err
	}

	err = db.WithContext(ctx).AutoMigrate(&User{})
	if err != nil {
		return err
	}

	err = db.WithContext(ctx).AutoMigrate(&UserFileProvider{})
	if err != nil {
		return err
	}

	err = db.WithContext(ctx).AutoMigrate(&GroupFileProvider{})
	if err != nil {
		return err
	}

	return nil
}
