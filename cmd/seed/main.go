package main

import (
	"log"

	"github.com/vyolayer/vyolayer/pkg/postgres"
	"github.com/vyolayer/vyolayer/pkg/postgres/seed"
)

func main() {
	// Load Config
	// cfg, err := config.Load("config/config.dev.yaml")
	// if err != nil {
	// 	panic(err)
	// }
	dsn := "postgres://vyolayer_user:vyolayer_password@localhost:4444/vyolayer_db?sslmode=disable"

	// Connect to DB
	db, err := postgres.NewConnectionFromDSN(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run Seeder
	if err := seed.Run(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("✅ Database seeded successfully!")
}
