package types

// Config represents the main configuration structure
type Config struct {
	Server   ServerConfig    `yaml:"server"`
	Services []ServiceConfig `yaml:"services"`
}
