package database

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, models ...interface{}) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if err := db.AutoMigrate(models...); err != nil {
		return err
	}

	return nil
}
