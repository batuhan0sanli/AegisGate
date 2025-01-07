package types

// ServerConfig holds server-related configurations
type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}
