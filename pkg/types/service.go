package types

// ServiceConfig holds configuration for a single service
type ServiceConfig struct {
	Name      string  `yaml:"name"`
	BasePath  string  `yaml:"base_path"`
	TargetURL string  `yaml:"target_url"`
	Routes    []Route `yaml:"routes"`
}

// Route represents a single route configuration
type Route struct {
	Path      string       `yaml:"path"`
	Methods   []HTTPMethod `yaml:"methods"`
	StripPath bool         `yaml:"strip_path"`
	Timeout   uint         `yaml:"timeout,omitempty"`
}

// expandMethods expands any abbreviations in the methods list and removes duplicates
func (r *Route) expandMethods() []HTTPMethod {
	methodSet := make(map[HTTPMethod]bool)
	expanded := make([]HTTPMethod, 0)

	// Expand all methods and add to set to remove duplicates
	for _, method := range r.Methods {
		for _, expandedMethod := range expandAbbreviation(method) {
			if !methodSet[expandedMethod] {
				methodSet[expandedMethod] = true
				expanded = append(expanded, expandedMethod)
			}
		}
	}

	return expanded
}

// GetMethods returns the expanded list of HTTP methods
func (r *Route) GetMethods() []HTTPMethod {
	return r.expandMethods()
}
