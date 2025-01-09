package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"AegisGate/internal/config"
	"AegisGate/internal/core"
)

// handleShutdown sets up signal handling and graceful shutdown
func handleShutdown(gateway *core.Gateway) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down gateway...", sig)
		
		if err := gateway.Close(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		os.Exit(0)
	}()
}

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

	// Create the gateway
	gateway, err := core.New(cfg, configPath)
	if err != nil {
		log.Fatalf("Failed to create gateway: %v", err)
	}
	defer gateway.Close()

	// Set up shutdown handling
	handleShutdown(gateway)

	// Start the gateway
	if err := gateway.Start(); err != nil {
		log.Fatalf("Gateway server failed: %v", err)
	}
}
