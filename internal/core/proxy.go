package core

import (
	"AegisGate/internal/logger"
	"AegisGate/pkg/types"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// ProxyManager manages reverse proxies for services
type ProxyManager struct {
	proxies map[string]*ServiceProxy
	mu      sync.RWMutex
	logger  *logger.Logger
}

// ServiceProxy represents a proxy configuration for a service
type ServiceProxy struct {
	name      string
	targetURL *url.URL
	proxy     *httputil.ReverseProxy
	config    types.ServiceConfig
	logger    *logger.RequestLogger
}

// NewProxyManager creates a new ProxyManager instance
func NewProxyManager(debug bool) *ProxyManager {
	return &ProxyManager{
		proxies: make(map[string]*ServiceProxy),
		logger:  logger.New(debug),
	}
}

// AddService creates and adds a new proxy for a service
func (pm *ProxyManager) AddService(service types.ServiceConfig) error {
	targetURL, err := url.Parse(service.TargetURL)
	if err != nil {
		return fmt.Errorf("invalid target URL for service %s: %w", service.Name, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	reqLogger := logger.NewRequestLogger(pm.logger, service.Name)

	// Configure proxy settings
	proxy.ModifyResponse = modifyResponse
	proxy.ErrorHandler = createErrorHandler(reqLogger)

	serviceProxy := &ServiceProxy{
		name:      service.Name,
		targetURL: targetURL,
		proxy:     proxy,
		config:    service,
		logger:    reqLogger,
	}

	pm.mu.Lock()
	pm.proxies[service.Name] = serviceProxy
	pm.mu.Unlock()

	return nil
}

// ServeHTTP handles the proxying of requests
func (sp *ServiceProxy) ServeHTTP(w http.ResponseWriter, r *http.Request, stripPath bool) {
	start := time.Now()

	// Log incoming request
	sp.logger.LogRequest(r)

	// Clone the request to modify it safely
	outReq := r.Clone(r.Context())

	// Strip path if configured
	if stripPath {
		originalPath := outReq.URL.Path
		outReq.URL.Path = stripBasePath(outReq.URL.Path, sp.config.BasePath)
		sp.logger.LogPathStripped(originalPath, outReq.URL.Path)
	}

	// Add custom headers
	outReq.Header.Set("X-Forwarded-Host", r.Host)
	outReq.Header.Set("X-Origin-Host", sp.targetURL.Host)
	outReq.Host = sp.targetURL.Host

	// Create a custom response writer to capture status code and size
	rw := logger.NewResponseWriter(w)

	// Forward the request to the target service
	sp.proxy.ServeHTTP(rw, outReq)

	// Log the completed request
	sp.logger.LogCompleted(r, rw, sp.targetURL.String()+outReq.URL.Path, start)
}

// createErrorHandler creates an error handler with logging
func createErrorHandler(reqLogger *logger.RequestLogger) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		reqLogger.LogError("Proxy error: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}
}

// GetProxy retrieves a proxy for a service
func (pm *ProxyManager) GetProxy(serviceName string) (*ServiceProxy, error) {
	pm.mu.RLock()
	proxy, exists := pm.proxies[serviceName]
	pm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("proxy not found for service: %s", serviceName)
	}

	return proxy, nil
}

// modifyResponse modifies the response from the backend service
func modifyResponse(resp *http.Response) error {
	// Add custom response headers
	resp.Header.Set("X-Proxy", "AegisGate")
	return nil
}

// stripBasePath removes the base path from the request path
func stripBasePath(path, basePath string) string {
	if basePath == "/" {
		return path
	}
	// Remove the base path while preserving the rest of the path
	return path[len(basePath):]
}
