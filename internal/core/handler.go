package core

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// handleNotFound returns a handler for 404 responses
func (g *Gateway) handleNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g.reqLogger.LogRequest(r)
		http.Error(w, "Not Found", http.StatusNotFound)
	})
}

// handleHealthCheck handles the health check endpoint
func (g *Gateway) handleHealthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	g.reqLogger.LogRequest(r)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
