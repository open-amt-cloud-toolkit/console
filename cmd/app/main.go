package main

import (
	"log"

	"github.com/open-amt-cloud-toolkit/console/config"
	"github.com/open-amt-cloud-toolkit/console/internal/app"
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
