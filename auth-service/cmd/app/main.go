package main

import (
	"log"

	"gitverse.ru/volatex/backend/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config errors: %s", err)
	}

}
