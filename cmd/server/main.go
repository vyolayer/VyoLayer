package main

import "github.com/vyolayer/vyolayer/internal/app"

func main() {
	server := app.New()
	server.LoadConfig()
	server.ConnectToDatabase()
	server.SetupMiddleware()
	server.SetupRoutes()

	if err := server.Run(); err != nil {
		panic(err)
	}

	server.ListenShutdownEvent()
}
