package main

import (
	"AegisGate/internal/logger"
	"AegisGate/internal/watcher"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"AegisGate/internal/config"
	"AegisGate/internal/core"
)

// handleShutdown sets up signal handling and graceful shutdown
func handleShutdown(g *core.Gateway, w *watcher.ConfigWatcher) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down gateway...", sig)

		// Close the config watcher
		if err := w.Close(); err != nil {
			log.Printf("Error closing config watcher: %v", err)
		}

		if err := g.Close(); err != nil {
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

	// Create the logger
	l := logger.New(cfg.Server.Debug)

	// Create the gateway
	gateway, err := core.New(cfg)
	if err != nil {
		log.Panicf("Failed to create gateway: %v", err)
	}
	defer func(gateway *core.Gateway) {
		log.Printf("Shutting down gateway...")
		err := gateway.Close()
		if err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}(gateway)

	// Initialize config watcher
	configWatcher, err := watcher.New(configPath, l)
	if err != nil {
		log.Fatalf("failed to create config watcher: %v", err)
	}
	defer func(configWatcher *watcher.ConfigWatcher) {
		err := configWatcher.Close()
		if err != nil {
			log.Printf("Error closing config watcher: %v", err)
		}
	}(configWatcher)
	configWatcher.RegisterHandler(gateway)

	// Set up shutdown handling
	handleShutdown(gateway, configWatcher)

	// Start the gateway
	err = gateway.Start()
	if !errors.Is(http.ErrServerClosed, err) {
		log.Fatalf("Gateway server failed: %v", err)
	}
}
