package model

import "gorm.io/gorm"

func MigrateModels(db *gorm.DB) error {
	err := db.AutoMigrate(&File{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&FileProvider{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&Group{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&UserGroup{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	return nil
}
