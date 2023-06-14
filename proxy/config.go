package proxy

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ListenerConfig holds the configuration for a listener, including the listener's address and the corresponding backend address for proxying requests.
type ListenerConfig struct {
	ListenerAddress  string   `yaml:"listenerAddress"`
	BackendAddresses []string `yaml:"backendAddresses"`
	TimeoutConnect   int      `yaml:"timeoutConnect"`
	TimeoutRead      int      `yaml:"timeoutRead"`
	TimeoutWrite     int      `yaml:"timeoutSend"`
}

// Config is a config structure for proxy application
type Config struct {
	Listeners []ListenerConfig `yaml:"listeners"`
	Debug     bool             `yaml:"debug"`
}

var (
	// ErrConfig is a global error for config-related errors.
	ErrConfig = errors.New("config error")
)

const (
	// DefaultConnectTimeout default timeout for establishing a connection (in seconds)
	DefaultConnectTimeout = 60
	// DefautlReadTimeout default timeout for read operations (in seconds)
	DefautlReadTimeout = 60
	// DefaultWriteTimeout default timeout for write operations (in seconds)
	DefaultWriteTimeout = 60
)

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

	for i, listener := range config.Listeners {
		if len(listener.BackendAddresses) == 0 {
			return nil, fmt.Errorf("%w: there is no any backends in %v listener", ErrConfig, listener.ListenerAddress)
		}

		if listener.TimeoutConnect == 0 {
			config.Listeners[i].TimeoutConnect = DefaultConnectTimeout
		}
		if listener.TimeoutRead == 0 {
			config.Listeners[i].TimeoutRead = DefautlReadTimeout
		}
		if listener.TimeoutWrite == 0 {
			config.Listeners[i].TimeoutWrite = DefaultWriteTimeout
		}
	}

	return &config, nil
}
