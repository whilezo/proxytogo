// Package proxy provides a TCP proxy server implementation for building proxy application.
package proxy

import (
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// Proxy is a proxy structure.
type Proxy struct {
	Backend        net.Conn
	Client         net.Conn
	TimeoutConnect time.Duration
	Debug          bool
}

// StartProxy starts the proxy server with the given configuration.
func StartProxy(listener *ListenerConfig, debug bool, wg *sync.WaitGroup) {
	defer wg.Done()

	server, err := net.Listen("tcp", listener.ListenerAddress)
	if err != nil {
		if debug {
			log.Printf("Error while starting server: %v", err)
		}
		return
	}
	if debug {
		log.Printf("Successfully started server on: %s", listener.ListenerAddress)
	}

	currentServerNum := 0
	globalProxy := new(Proxy)
	globalProxy.Debug = debug
	globalProxy.TimeoutConnect = time.Duration(listener.TimeoutConnect) * time.Second

	for {
		conn, err := server.Accept()
		if err != nil {
			if debug {
				log.Printf("Error while accepting client: %v", err)
			}
			continue
		}
		if debug {
			log.Printf("New client: %s", conn.RemoteAddr().String())
		}
		proxy := *globalProxy
		proxy.Client = conn

		backend, err := net.DialTimeout("tcp", listener.BackendAddresses[currentServerNum], proxy.TimeoutConnect)
		if err != nil {
			if debug {
				log.Printf("Error while connecting to backend: %v", err)
			}
			continue
		}
		if debug {
			log.Printf("Connected to backend: %s", backend.RemoteAddr().String())
		}
		proxy.Backend = backend
		currentServerNum = (currentServerNum + 1) % len(listener.BackendAddresses)

		go proxy.ForwardRequests()
	}
}

// ForwardRequests recieves request from client and forwards it to proxy and backwards.
func (p *Proxy) ForwardRequests() {
	buffer := make([]byte, 4096)

	close := func() {
		p.Client.Close()
		p.Backend.Close()
	}

	// Reading from client and writing to backend.
	go func() {
		defer close()

		n, err := io.CopyBuffer(p.Backend, p.Client, buffer)
		if p.Debug {
			log.Printf("Incoming TCP connection closed; error: %v; bytes forwarded: %d\n", err, n)
		}
	}()

	// Reading from backend and writing to client.
	go func() {
		defer close()

		n, err := io.CopyBuffer(p.Client, p.Backend, buffer)
		if p.Debug {
			log.Printf("Outgoing TCP connection closed; error: %v; bytes forwarded: %d\n", err, n)
		}
	}()
}
