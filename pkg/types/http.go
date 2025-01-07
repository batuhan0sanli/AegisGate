package types

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

// HTTPMethod represents an HTTP method
type HTTPMethod string

// HTTP methods as defined in RFC 7231 and RFC 5789
const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	PATCH   HTTPMethod = "PATCH"
	HEAD    HTTPMethod = "HEAD"
	OPTIONS HTTPMethod = "OPTIONS"
	TRACE   HTTPMethod = "TRACE"
	CONNECT HTTPMethod = "CONNECT"
)

// IsValid checks if the HTTP method is valid
func (m *HTTPMethod) IsValid() bool {
	switch *m {
	case GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE, CONNECT:
		return true
	default:
		return false
	}
}

// String returns the string representation of the HTTP method
func (m *HTTPMethod) String() string {
	return string(*m)
}

// ParseHTTPMethod converts a string to an HTTPMethod and validates it
func ParseHTTPMethod(s string) (*HTTPMethod, error) {
	method := HTTPMethod(strings.ToUpper(s))
	if !method.IsValid() {
		return nil, fmt.Errorf("invalid HTTP method: %s", s)
	}
	return &method, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (m *HTTPMethod) UnmarshalYAML(value *yaml.Node) error {
	method := HTTPMethod(strings.ToUpper(value.Value))
	if !method.IsValid() {
		return fmt.Errorf("invalid HTTP method: %s", value.Value)
	}
	*m = method
	return nil
}
