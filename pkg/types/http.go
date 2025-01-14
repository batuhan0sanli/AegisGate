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

// Abbreviated HTTP methods
const (
	CRUD HTTPMethod = "CRUD" // Create, Read, Update, Delete
	RO   HTTPMethod = "RO"   // Read-Only (GET, HEAD)
	RW   HTTPMethod = "RW"   // Read-Write (GET, POST, PUT, PATCH)
	FULL HTTPMethod = "FULL" // All common methods including OPTIONS, TRACE and CONNECT
)

// expandAbbreviation converts method abbreviations to their corresponding HTTP methods
func expandAbbreviation(method HTTPMethod) []HTTPMethod {
	switch method {
	case CRUD:
		return []HTTPMethod{GET, POST, PUT, PATCH, DELETE}
	case RO:
		return []HTTPMethod{GET, HEAD}
	case RW:
		return []HTTPMethod{GET, POST, PUT, PATCH}
	case FULL:
		return []HTTPMethod{GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE, CONNECT}
	default:
		return []HTTPMethod{method}
	}
}

// IsValid checks if the HTTP method is valid
func (m *HTTPMethod) IsValid() bool {
	switch *m {
	case GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE, CONNECT,
		CRUD, RO, RW, FULL:
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
	method, err := ParseHTTPMethod(value.Value)
	if err != nil {
		return err
	}
	*m = *method
	return nil
}
