package main

import (
	"gitverse.ru/volatex/backend/internal/app"
	"log"

	"gitverse.ru/volatex/backend/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config errors: %s", err)
	}

	app.Run(cfg)
}
