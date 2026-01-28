package main

import (
	"worklayer/internal/config"
	"worklayer/internal/platform/database"
	"worklayer/internal/platform/database/models"
)

func main() {
	cfg, err := config.Load("config/config.dev.yaml")
	if err != nil {
		panic(err)
	}

	db, err := database.Init(&cfg.Database)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(
		db,
		models.User{},
		models.UserSession{},
	)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	defer sqlDB.Close()
}
