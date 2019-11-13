package main

import (
	"github.com/AdrianOrlow/files-api/app"
	"github.com/AdrianOrlow/files-api/config"
	"log"
)

func main() {
	config, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	app := &app.App{}
	app.Initialize(config)
	app.Run(":8000")
}