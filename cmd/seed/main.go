package main

import (
	"log"
	"worklayer/internal/config"
	"worklayer/internal/platform/database"
	"worklayer/internal/platform/database/seed"
)

func main() {
	// Load Config
	cfg, err := config.Load("config/config.dev.yaml")
	if err != nil {
		panic(err)
	}

	// Connect to DB
	db, err := database.Init(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run Seeder
	if err := seed.Run(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("✅ Database seeded successfully!")
}
