// Package proxy provides a TCP proxy server implementation for building proxy application.
package proxy

import (
	"io"
	"log"
	"net"
)

// Proxy is a proxy structure.
type Proxy struct {
	Backend net.Conn
	Client  net.Conn
	Debug   bool
}

// StartProxy starts the proxy server with the given configuration.
func StartProxy(listener *ListenerConfig, debug bool) {
	server, err := net.Listen("tcp", listener.ListenerAddress)
	if err != nil {
		if debug {
			log.Printf("Error while starting server: %v", err)
		}
		return
	}
	if debug {
		log.Print("Successfully started server")
	}

	currentServerNum := 0

	for {
		proxy := new(Proxy)
		proxy.Debug = debug

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
		proxy.Client = conn

		backend, err := net.Dial("tcp", listener.BackendAddresses[currentServerNum])
		if err != nil {
			if debug {
				log.Printf("Error while connecting to backend: %v", err)
			}
			continue
		}
		if debug {
			log.Printf("Connected to backend: %s", backend.RemoteAddr().String())
		}
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
