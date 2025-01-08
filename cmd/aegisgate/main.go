package main

import (
	"log"
	"os"

	"AegisGate/internal/config"
	"AegisGate/internal/core"
)

func main() {
	// Get config file path from command line argument or use default
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	// Load the configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start the gateway
	gateway := core.New(cfg)
	if err := gateway.Start(); err != nil {
		log.Fatalf("Gateway server failed: %v", err)
	}
}
