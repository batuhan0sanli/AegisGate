package watcher

import (
	"AegisGate/internal/config"
	"AegisGate/internal/logger"
	"AegisGate/pkg/types"
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// ConfigWatcher watches for configuration file changes
type ConfigWatcher struct {
	watcher    *fsnotify.Watcher
	logger     *logger.Logger
	configPath string
	mu         sync.RWMutex
	handlers   []ConfigChangeHandler
}

// ConfigChangeHandler is called when configuration changes
type ConfigChangeHandler interface {
	OnConfigChange(newConfig *types.Config) error
}

// New creates a new ConfigWatcher
func New(configPath string, logger *logger.Logger) (*ConfigWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	cw := &ConfigWatcher{
		watcher:    watcher,
		logger:     logger,
		configPath: configPath,
		handlers:   make([]ConfigChangeHandler, 0),
	}

	if err := cw.watch(); err != nil {
		return nil, err
	}

	return cw, nil
}

// RegisterHandler adds a new configuration change handler
func (cw *ConfigWatcher) RegisterHandler(handler ConfigChangeHandler) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.handlers = append(cw.handlers, handler)
}

// watch starts watching the configuration file
func (cw *ConfigWatcher) watch() error {
	if err := cw.watcher.Add(cw.configPath); err != nil {
		return fmt.Errorf("failed to watch config file: %w", err)
	}

	go cw.watchLoop()
	return nil
}

// watchLoop handles file system events
func (cw *ConfigWatcher) watchLoop() {
	for {
		select {
		case event, ok := <-cw.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				cw.handleConfigChange()
			}
		case err, ok := <-cw.watcher.Errors:
			if !ok {
				return
			}
			cw.logger.Error("Config watcher error: %v", err)
		}
	}
}

// handleConfigChange processes configuration changes
func (cw *ConfigWatcher) handleConfigChange() {
	cw.logger.Info("Configuration file changed, reloading...")

	// Load new configuration
	newConfig, err := config.LoadConfig(cw.configPath)
	if err != nil {
		cw.logger.Error("Failed to load new configuration: %v", err)
		return
	}

	// Notify all handlers
	cw.mu.RLock()
	defer cw.mu.RUnlock()

	for _, handler := range cw.handlers {
		if err := handler.OnConfigChange(newConfig); err != nil {
			cw.logger.Error("Handler failed to process config change: %v", err)
		}
	}

	cw.logger.Info("Configuration reloaded successfully")
}

// Close stops watching and cleans up resources
func (cw *ConfigWatcher) Close() error {
	cw.logger.Debug("Closing config watcher")
	return cw.watcher.Close()
}
