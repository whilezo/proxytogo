package proxy

import (
	"net"
	"time"
)

// HealthCheck holds configuration for health checks.
type HealthCheck struct {
	BackendAddr   string
	Status        bool
	CheckInterval time.Duration
	Timeout       time.Duration
}

// GetAvailableBackend gets first available backend server in map and returns its adress as string.
func GetAvailableBackend(healthStatus map[string]*bool) string {
	for address, status := range healthStatus {
		if *status {
			return address
		}
	}

	return ""
}

// StartHealthCheck runs a health check for all listener backends
func StartHealthCheck(backends []string, interval, timeout time.Duration) map[string]*bool {
	backendStatus := make(map[string]*bool, 0)

	for _, backend := range backends {
		h := HealthCheck{
			BackendAddr:   backend,
			Status:        false,
			CheckInterval: interval,
			Timeout:       timeout,
		}
		go h.Run()
		backendStatus[backend] = &h.Status
	}

	return backendStatus
}

// Run starts to perform a health check on the backend.
func (h *HealthCheck) Run() {
	for {
		conn, err := net.DialTimeout("tcp", h.BackendAddr, h.Timeout)

		if err != nil {
			h.Status = false
		} else {
			conn.Close()
			h.Status = true
		}

		time.Sleep(h.CheckInterval)
	}
}
