package main

import (
	"github.com/joho/godotenv"
	"github.com/quinntas/go-fiber-template/database"
	"github.com/quinntas/go-fiber-template/eventEmitter"
	"github.com/quinntas/go-fiber-template/resources"
	"github.com/quinntas/go-fiber-template/server"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.SetupDatabase()
	eventEmitter.SetupEventEmitter()

	app := server.Create()

	resources.SetupRouter(app)
	resources.SetupEvents(eventEmitter.Manager)

	if err := server.Listen(app); err != nil {
		log.Fatal(err)
	}
}
