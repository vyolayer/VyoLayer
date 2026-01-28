package main

import (
	"log"

	"worklayer/internal/app"
)

func main() {
	server := app.New()
	server.LoadConfig()
	server.ConnectToDatabase()
	server.SetupMiddleware()
	server.SetupRoutes()

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

	server.ListenShutdownEvent()
}
