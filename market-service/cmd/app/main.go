package main

import (
	"gitverse.ru/volatex/backend/market-service/config"
	"gitverse.ru/volatex/backend/market-service/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
