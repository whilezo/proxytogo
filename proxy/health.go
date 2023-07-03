package proxy

import (
	"fmt"
	"net"
	"time"
)

// HealthCheck holds configuration for health checks.
type HealthCheck struct {
	BackendAddr   string
	Status        bool
	CheckInterval time.Duration
	Timeout       time.Duration
	Debug         bool
}

// StartHealthCheck runs a health check for all listener backends
func StartHealthCheck(backends []string, interval, timeout int, debug bool) map[string]*bool {
	backendStatus := make(map[string]*bool, 0)

	for _, backend := range backends {
		h := HealthCheck{
			BackendAddr:   backend,
			Status:        false,
			CheckInterval: time.Duration(interval) * time.Second,
			Timeout:       time.Duration(timeout) * time.Second,
			Debug:         debug,
		}
		go h.Run()
		backendStatus[backend] = &h.Status
	}

	return backendStatus
}

// Run starts to perform a health check on the backend.
func (h *HealthCheck) Run() {
	for {
		_, err := net.DialTimeout("tcp", h.BackendAddr, h.Timeout)

		if err != nil {
			fmt.Println(err)
			h.Status = false
		} else {
			h.Status = true
		}

		time.Sleep(h.CheckInterval)
	}
}
