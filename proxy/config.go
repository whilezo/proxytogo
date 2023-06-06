package proxy

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ListenerConfig holds the configuration for a listener, including the listener's address and the corresponding backend address for proxying requests.
type ListenerConfig struct {
	ListenerAddress  string   `yaml:"listenerAddress"`
	BackendAddresses []string `yaml:"backendAddresses"`
}

// Config is a config structure for proxy application
type Config struct {
	Listeners []ListenerConfig `yaml:"listeners"`
	Debug     bool             `yaml:"debug"`
}

// LoadConfig loads the YAML configuration file and returns the parsed configuration data.
func LoadConfig(path string) (*Config, error) {
	var config Config

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
