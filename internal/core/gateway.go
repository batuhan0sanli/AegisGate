package core

import (
	"AegisGate/internal/logger"
	"AegisGate/pkg/types"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Gateway represents the API gateway
type Gateway struct {
	config    *types.Config
	router    *httprouter.Router
	proxies   *ProxyManager
	logger    *logger.Logger
	reqLogger *logger.RequestLogger
}

// New creates a new Gateway instance
func New(config *types.Config, configPath string) (*Gateway, error) {
	l := logger.New(config.Server.Debug)

	g := &Gateway{
		config:    config,
		router:    httprouter.New(),
		proxies:   NewProxyManager(config.Server.Debug),
		logger:    l,
		reqLogger: logger.NewRequestLogger(l, "AegisGate"),
	}

	g.logger.Debug("Debug mode enabled")

	// Set up 404 handler
	g.router.NotFound = g.handleNotFound()

	// Initialize routes
	if err := g.initializeRoutes(); err != nil {
		return nil, fmt.Errorf("failed to initialize routes: %w", err)
	}

	return g, nil
}

// initializeRoutes sets up all the routes from the configuration
func (g *Gateway) initializeRoutes() error {
	for _, service := range g.config.Services {
		// Add service to proxy manager
		if err := g.proxies.AddService(service); err != nil {
			return fmt.Errorf("failed to add service proxy: %w", err)
		}

		// Set up routes for the service
		for _, route := range service.Routes {
			routerPath := g.convertPath(service.BasePath, route.Path)

			// Register handlers for each HTTP method
			for _, method := range route.Methods {
				handler := g.createHandler(service, route)
				g.router.Handle(method.String(), routerPath, handler)
				g.logger.Debug("Registered route: %s %s -> %s", method, routerPath, service.TargetURL)
			}
		}
	}

	return nil
}

// convertPath combines base path and route path, ensuring proper formatting
func (g *Gateway) convertPath(basePath, routePath string) string {
	// Combine paths and ensure single forward slash between segments
	fullPath := fmt.Sprintf("%s/%s", strings.TrimSuffix(basePath, "/"), strings.TrimPrefix(routePath, "/"))

	// Convert path parameters from {param} to :param format for httprouter
	segments := strings.Split(fullPath, "/")
	for i, segment := range segments {
		switch {
		// Handle named parameters {param}
		case strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}"):
			segments[i] = ":" + segment[1:len(segment)-1]
		// Handle wildcards *
		case segment == "*":
			segments[i] = "*path" // httprouter requires named wildcards
		}
	}

	return strings.Join(segments, "/")
}

// createHandler creates a handler function for a specific route
func (g *Gateway) createHandler(service types.ServiceConfig, route types.Route) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the proxy for this service
		proxy, err := g.proxies.GetProxy(service.Name)
		if err != nil {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}

		// Apply timeout if specified
		if route.Timeout > 0 {
			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(route.Timeout)*time.Second)
			defer cancel()
			r = r.WithContext(ctx)
		}

		// Copy URL parameters to request
		for _, p := range ps {
			r.URL.Query().Set(p.Key, p.Value)
		}

		// Forward the request to the target service
		proxy.ServeHTTP(w, r, route.StripPath)
	}
}

// Start starts the gateway server
func (g *Gateway) Start() error {
	addr := fmt.Sprintf("%s:%d", g.config.Server.Host, g.config.Server.Port)
	g.logger.Info("Starting gateway server on %s", addr)
	return http.ListenAndServe(addr, g.router)
}

// handleNotFound returns a handler for 404 responses
func (g *Gateway) handleNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g.reqLogger.LogRequest(r)
		http.Error(w, "Not Found", http.StatusNotFound)
	})
}

// Add Close method
func (g *Gateway) Close() error {
	if g.watcher != nil {
		return g.watcher.Close()
	}
	return nil
}
