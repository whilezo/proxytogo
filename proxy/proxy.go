// Package proxy provides a TCP proxy server implementation for building proxy application.
package proxy

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Proxy is a proxy configuration for backend and client.
type Proxy struct {
	Backend          net.Conn
	BackendAddresses []string
	Client           net.Conn
	Health           map[string]*bool
	Protocol         string
	TimeoutConnect   time.Duration
	TimeoutRead      time.Duration
	TimeoutWrite     time.Duration
	Debug            bool
}

// StartProxy starts the proxy server with the given configuration.
func StartProxy(listener *ListenerConfig, debug bool, wg *sync.WaitGroup) error {
	var TCPServer net.Listener
	var UDPServer *net.UDPConn
	var err error

	defer wg.Done()

	if listener.Protocol == "tcp" {
		TCPServer, err = net.Listen("tcp", listener.ListenerAddress)
		if err != nil {
			if debug {
				log.Printf("Error while starting server: %v", err)
			}
			return err
		}
	} else if listener.Protocol == "udp" {
		UDPAddr, err := net.ResolveUDPAddr("udp", listener.ListenerAddress)
		if err != nil {
			return err
		}

		UDPServer, err = net.ListenUDP("udp", UDPAddr)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported protocol")
	}

	if debug {
		log.Printf("Successfully started server on: %s", listener.ListenerAddress)
	}

	globalProxy := new(Proxy)
	globalProxy.BackendAddresses = listener.BackendAddresses
	globalProxy.Health = StartHealthCheck(
		listener.BackendAddresses,
		time.Duration(listener.HealthCheckInterval)*time.Second,
		time.Duration(listener.TimeoutConnect)*time.Second,
		debug,
	)
	globalProxy.Protocol = listener.Protocol
	globalProxy.TimeoutConnect = time.Duration(listener.TimeoutConnect) * time.Second
	globalProxy.TimeoutRead = time.Duration(listener.TimeoutRead) * time.Second
	globalProxy.TimeoutWrite = time.Duration(listener.TimeoutWrite) * time.Second
	globalProxy.Debug = debug

	if globalProxy.Protocol == "tcp" {
		globalProxy.listenTCP(TCPServer, listener)
	}

	if globalProxy.Protocol == "udp" {
		globalProxy.listenUDP(UDPServer)
	}

	return nil
}

func (p Proxy) listenUDP(s *net.UDPConn) {
	currentServerNum := 0

	for {
		backendAddr := p.BackendAddresses[currentServerNum]

		if !*p.Health[backendAddr] {
			if p.Debug {
				log.Printf("Backend %s doesn't work", backendAddr)
			}

			backendAddr = GetAvailableBackend(p.Health)
			if backendAddr == "" {
				log.Printf("No available backends right now")
				continue
			}
		}

		backend, err := net.DialTimeout("tcp", backendAddr, p.TimeoutConnect)
		if err != nil {
			if p.Debug {
				log.Printf("Error while connecting to backend: %v", err)
			}
			continue
		}
		if p.Debug {
			log.Printf("Connected to backend: %s", backend.RemoteAddr().String())
		}
		p.Backend = backend
		currentServerNum = (currentServerNum + 1) % len(p.BackendAddresses)

		go func() {
			var err error
			var bytesForwarded int

			for {
				buffer := make([]byte, 4096)

				var n int

				n, _, err = s.ReadFromUDP(buffer)
				if err != nil {
					break
				}

				n, err = p.Backend.Write(buffer)
				if err != nil {
					break
				}

				bytesForwarded += n
			}

			if p.Debug {
				log.Printf("Incoming UDP connection closed; error: %v; bytes forwarded: %d\n", err, bytesForwarded)
			}
		}()

		go func() {
			var err error
			var bytesForwarded int

			for {
				buffer := make([]byte, 4096)

				var n int

				n, err = p.Backend.Read(buffer)
				if err != nil {
					break
				}

				n, _, err = s.ReadFromUDP(buffer)
				if err != nil {
					break
				}

				bytesForwarded += n
			}

			if p.Debug {
				log.Printf("Incoming UDP connection closed; error: %v; bytes forwarded: %d\n", err, bytesForwarded)
			}
		}()
	}
}

func (p Proxy) listenTCP(s net.Listener, cfg *ListenerConfig) {
	currentServerNum := 0

	for {
		conn, err := s.Accept()
		if err != nil {
			if p.Debug {
				log.Printf("Error while accepting client: %v", err)
			}
			continue
		}
		if p.Debug {
			log.Printf("New client: %s", conn.RemoteAddr().String())
		}

		p.Client = conn

		backendAddr := p.BackendAddresses[currentServerNum]

		if !*p.Health[backendAddr] {
			if p.Debug {
				log.Printf("Backend %s doesn't work", backendAddr)
			}

			backendAddr = GetAvailableBackend(p.Health)
			if backendAddr == "" {
				log.Printf("No available backends right now")
				continue
			}
		}

		backend, err := net.DialTimeout("tcp", backendAddr, p.TimeoutConnect)
		if err != nil {
			if p.Debug {
				log.Printf("Error while connecting to backend: %v", err)
			}
			conn.Close()
			continue
		}
		if p.Debug {
			log.Printf("Connected to backend: %s", backend.RemoteAddr().String())
		}
		p.Backend = backend
		currentServerNum = (currentServerNum + 1) % len(cfg.BackendAddresses)

		go p.ForwardRequests()
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
