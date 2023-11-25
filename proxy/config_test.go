package proxy_test

import (
	"proxy/proxy"
	"testing"
)

func TestConfig(t *testing.T) {
	// Test case 1: Valid configuration
	config := proxy.Config{
		Listeners: []proxy.ListenerConfig{
			{
				ListenerAddress:  "localhost:8000",
				BackendAddresses: []string{"backend1:9000", "backend2:9001"},
				TimeoutConnect:   5000,
			},
		},
		Debug: true,
	}

	// Assert the values of the configuration
	if config.Debug != true {
		t.Error("Expected Debug to be true, but got false")
	}

	if len(config.Listeners) != 1 {
		t.Errorf("Expected 1 listener, but got %d", len(config.Listeners))
	}

	listener := config.Listeners[0]
	if listener.ListenerAddress != "localhost:8000" {
		t.Errorf("Expected ListenerAddress to be localhost:8000, but got %s", listener.ListenerAddress)
	}

	if len(listener.BackendAddresses) != 2 {
		t.Errorf("Expected 2 backend addresses, but got %d", len(listener.BackendAddresses))
	}

	if listener.TimeoutConnect != 5000 {
		t.Errorf("Expected timeoutConnect to be 5000, but got %d", listener.TimeoutConnect)
	}

	// Test case 2: Empty configuration
	emptyConfig := proxy.Config{}
	if emptyConfig.Debug != false {
		t.Error("Expected Debug to be false, but got true")
	}

	if len(emptyConfig.Listeners) != 0 {
		t.Errorf("Expected 0 listeners, but got %d", len(emptyConfig.Listeners))
	}
}
