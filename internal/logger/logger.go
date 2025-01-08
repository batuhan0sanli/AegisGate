package logger

import (
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

// Logger represents the logging configuration
type Logger struct {
	debug bool
}

// New creates a new Logger instance
func New(debug bool) *Logger {
	return &Logger{
		debug: debug,
	}
}

// Debug logs debug messages when debug mode is enabled
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.debug {
		log.Printf("[DEBUG] "+format, v...)
	}
}

// ServiceDebug logs debug messages with service name
func (l *Logger) ServiceDebug(serviceName, format string, v ...interface{}) {
	if l.debug {
		log.Printf("[DEBUG] [%s] "+format, append([]interface{}{serviceName}, v...)...)
	}
}

// Error logs error messages
func (l *Logger) Error(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

// Info logs info messages
func (l *Logger) Info(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// ResponseWriter wraps http.ResponseWriter to capture status code and response size
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// WriteHeader captures the status code and calls the underlying WriteHeader
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the response size and calls the underlying Write
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}

// RequestLogger handles request logging
type RequestLogger struct {
	logger      *Logger
	serviceName string
}

// NewRequestLogger creates a new RequestLogger
func NewRequestLogger(logger *Logger, serviceName string) *RequestLogger {
	return &RequestLogger{
		logger:      logger,
		serviceName: serviceName,
	}
}

// LogRequest logs the incoming request details
func (rl *RequestLogger) LogRequest(r *http.Request) {
	rl.logger.ServiceDebug(rl.serviceName, "Incoming request: %s %s", r.Method, r.URL.Path)
	if rl.logger.debug {
		if dump, err := httputil.DumpRequest(r, true); err == nil {
			rl.logger.ServiceDebug(rl.serviceName, "Request details:\n%s", string(dump))
		}
	}
}

// LogPathStripped logs path stripping information
func (rl *RequestLogger) LogPathStripped(originalPath, newPath string) {
	rl.logger.ServiceDebug(rl.serviceName, "Path stripped: %s -> %s", originalPath, newPath)
}

// LogCompleted logs the completed request details
func (rl *RequestLogger) LogCompleted(r *http.Request, rw *ResponseWriter, targetURL string, start time.Time) {
	duration := time.Since(start)
	rl.logger.ServiceDebug(rl.serviceName,
		"Completed %s %s -> %s [%d] (%d bytes) in %v",
		r.Method,
		r.URL.Path,
		targetURL,
		rw.statusCode,
		rw.size,
		duration,
	)
}

// LogError logs error messages with service context
func (rl *RequestLogger) LogError(format string, v ...interface{}) {
	rl.logger.Error("[%s] "+format, append([]interface{}{rl.serviceName}, v...)...)
}
