package core

import (
	"AegisGate/pkg/types"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

// ProxyManager manages reverse proxies for services
type ProxyManager struct {
	proxies map[string]*ServiceProxy
	mu      sync.RWMutex
}

// ServiceProxy represents a proxy configuration for a service
type ServiceProxy struct {
	name      string
	targetURL *url.URL
	proxy     *httputil.ReverseProxy
	config    types.ServiceConfig
}

// NewProxyManager creates a new ProxyManager instance
func NewProxyManager() *ProxyManager {
	return &ProxyManager{
		proxies: make(map[string]*ServiceProxy),
	}
}

// AddService creates and adds a new proxy for a service
func (pm *ProxyManager) AddService(service types.ServiceConfig) error {
	targetURL, err := url.Parse(service.TargetURL)
	if err != nil {
		return fmt.Errorf("invalid target URL for service %s: %w", service.Name, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Configure proxy settings
	proxy.ModifyResponse = modifyResponse
	proxy.ErrorHandler = errorHandler

	serviceProxy := &ServiceProxy{
		name:      service.Name,
		targetURL: targetURL,
		proxy:     proxy,
		config:    service,
	}

	pm.mu.Lock()
	pm.proxies[service.Name] = serviceProxy
	pm.mu.Unlock()

	return nil
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

// ServeHTTP handles the proxying of requests
func (sp *ServiceProxy) ServeHTTP(w http.ResponseWriter, r *http.Request, stripPath bool) {
	// Clone the request to modify it safely
	outReq := r.Clone(r.Context())

	// Strip path if configured
	if stripPath {
		outReq.URL.Path = stripBasePath(outReq.URL.Path, sp.config.BasePath)
	}

	// Add custom headers
	outReq.Header.Set("X-Forwarded-Host", r.Host)
	outReq.Header.Set("X-Origin-Host", sp.targetURL.Host)
	outReq.Host = sp.targetURL.Host

	sp.proxy.ServeHTTP(w, outReq)
}

// modifyResponse modifies the response from the backend service
func modifyResponse(resp *http.Response) error {
	// Add custom response headers
	resp.Header.Set("X-Proxy", "AegisGate")
	return nil
}

// errorHandler handles proxy errors
func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	// Log the error
	fmt.Printf("Proxy error: %v\n", err)

	// Return a 502 Bad Gateway error
	http.Error(w, "Bad Gateway", http.StatusBadGateway)
}

// stripBasePath removes the base path from the request path
func stripBasePath(path, basePath string) string {
	if basePath == "/" {
		return path
	}
	// Remove the base path while preserving the rest of the path
	return path[len(basePath):]
} 