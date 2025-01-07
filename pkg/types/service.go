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
