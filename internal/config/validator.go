package config

import (
	"AegisGate/pkg/types"
	"fmt"
	"net/url"
	"strings"
)

// validateConfig performs basic validation of the configuration
func validateConfig(config *types.Config) error {
	if err := validateServer(config.Server); err != nil {
		return fmt.Errorf("server validation failed: %w", err)
	}

	if err := validateServices(config.Services); err != nil {
		return fmt.Errorf("services validation failed: %w", err)
	}

	return nil
}

// validateServer validates server-specific configuration
func validateServer(server types.ServerConfig) error {
	if server.Port <= 0 || server.Port > 65535 {
		return fmt.Errorf("invalid port number: %d (must be between 1 and 65535)", server.Port)
	}

	if server.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	return nil
}

// validateServices validates the services configuration
func validateServices(services []types.ServiceConfig) error {
	if len(services) == 0 {
		return fmt.Errorf("at least one service must be configured")
	}

	serviceNames := make(map[string]bool)
	servicePaths := make(map[string]bool)

	for i, service := range services {
		if err := validateService(service, i); err != nil {
			return err
		}

		if serviceNames[service.Name] {
			return fmt.Errorf("service[%d]: duplicate service name '%s'", i, service.Name)
		}
		serviceNames[service.Name] = true

		if servicePaths[service.BasePath] {
			return fmt.Errorf("service[%d]: duplicate base path '%s'", i, service.BasePath)
		}
		servicePaths[service.BasePath] = true
	}

	return nil
}

// validateService validates a single service configuration
func validateService(service types.ServiceConfig, index int) error {
	reservedPaths := map[string]bool{
		"/health": true,
	}

	if service.Name == "" {
		return fmt.Errorf("service[%d]: name cannot be empty", index)
	}

	if service.BasePath == "" {
		return fmt.Errorf("service[%d]: base path cannot be empty", index)
	}

	if reservedPaths[service.BasePath] {
		return fmt.Errorf("service[%s]: path '%s' is reserved for internal use", service.Name, service.BasePath)
	}

	if !strings.HasPrefix(service.BasePath, "/") {
		return fmt.Errorf("service[%d]: base path must start with '/'", index)
	}

	if service.TargetURL == "" {
		return fmt.Errorf("service[%d]: target URL cannot be empty", index)
	}

	if _, err := url.Parse(service.TargetURL); err != nil {
		return fmt.Errorf("service[%d]: invalid target URL '%s': %v", index, service.TargetURL, err)
	}

	if err := validateRoutes(service.Routes, index); err != nil {
		return err
	}

	return nil
}

// validateRoutes validates the routes configuration for a service
func validateRoutes(routes []types.Route, serviceIndex int) error {
	if len(routes) == 0 {
		return fmt.Errorf("service[%d]: at least one route must be configured", serviceIndex)
	}

	routePaths := make(map[string]bool)

	for i, route := range routes {
		if err := validateRoute(route, serviceIndex, i); err != nil {
			return err
		}

		if routePaths[route.Path] {
			return fmt.Errorf("service[%d].route[%d]: duplicate path '%s'", serviceIndex, i, route.Path)
		}
		routePaths[route.Path] = true
	}

	return nil
}

// validateRoute validates a single route configuration
func validateRoute(route types.Route, serviceIndex, routeIndex int) error {
	if route.Path == "" {
		return fmt.Errorf("service[%d].route[%d]: path cannot be empty", serviceIndex, routeIndex)
	}

	if !strings.HasPrefix(route.Path, "/") {
		return fmt.Errorf("service[%d].route[%d]: path must start with '/'", serviceIndex, routeIndex)
	}

	if len(route.Methods) == 0 {
		return fmt.Errorf("service[%d].route[%d]: at least one HTTP method must be specified", serviceIndex, routeIndex)
	}

	// Validate HTTP methods
	for _, method := range route.Methods {
		if !method.IsValid() {
			return fmt.Errorf("service[%d].route[%d]: invalid HTTP method '%s'", serviceIndex, routeIndex, method.String())
		}
	}

	// Validate timeout format if specified
	if route.Timeout != 0 {
		if err := validateTimeout(route.Timeout, serviceIndex, routeIndex); err != nil {
			return err
		}
	}

	return nil
}

// validateTimeout validates the timeout format
func validateTimeout(timeout uint, serviceIndex, routeIndex int) error {
	if timeout <= 0 {
		return fmt.Errorf("service[%d].route[%d]: timeout must be greater than 0", serviceIndex, routeIndex)
	}

	return nil
}
