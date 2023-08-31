package main

import (
	"avito-segments/config"
	"avito-segments/internal/app"
	"log"
)

func main() {
	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
