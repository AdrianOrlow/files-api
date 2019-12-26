package main

import (
	"github.com/AdrianOrlow/files-api/app"
	"github.com/AdrianOrlow/files-api/config"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	config := config.LoadConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":8000")
}
