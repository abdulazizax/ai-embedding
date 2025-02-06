package main

import (
	"log"

	"github.com/abdulazizax/ai-embedding/config"
	"github.com/abdulazizax/ai-embedding/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
