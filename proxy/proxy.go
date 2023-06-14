// Package proxy provides a TCP proxy server implementation for building proxy application.
package proxy

import (
	"errors"
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
	TimeoutRead    time.Duration
	TimeoutWrite   time.Duration
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
	globalProxy.TimeoutConnect = time.Duration(listener.TimeoutConnect) * time.Second
	globalProxy.TimeoutRead = time.Duration(listener.TimeoutRead) * time.Second
	globalProxy.TimeoutWrite = time.Duration(listener.TimeoutWrite) * time.Second
	globalProxy.Debug = debug

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
			conn.Close()
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
	closeConnections := func() {
		p.Client.Close()
		p.Backend.Close()
	}

	// Reading from client and writing to backend.
	go func() {
		defer closeConnections()
		var n int
		var err error
		var bytesForwarded int
		buffer := make([]byte, 4096)

		for {
			p.Client.SetReadDeadline(time.Now().Add(p.TimeoutRead))
			p.Client.SetWriteDeadline(time.Now().Add(p.TimeoutWrite))

			n, err = copyBuffer(p.Client, p.Backend, buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					err = nil
				}
				break
			}
			bytesForwarded += n
		}

		if p.Debug {
			log.Printf("Incoming TCP connection closed; error: %v; bytes forwarded: %d\n", err, bytesForwarded)
		}
	}()

	// Reading from backend and writing to client.
	go func() {
		defer closeConnections()
		var n int
		var err error
		var bytesForwarded int
		buffer := make([]byte, 4096)

		for {
			n, err = copyBuffer(p.Backend, p.Client, buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					err = nil
				}
				break
			}
			bytesForwarded += n
		}

		if p.Debug {
			log.Printf("Outgoing TCP connection closed; error: %v; bytes forwarded: %d\n", err, bytesForwarded)
		}
	}()
}

func copyBuffer(src, dst net.Conn, buffer []byte) (int, error) {
	n, err := src.Read(buffer)
	if err != nil {
		return 0, err
	}

	n, err = dst.Write(buffer[:n])
	if err != nil {
		return n, err
	}

	return n, nil
}
