package udp

import (
	"fmt"
	"net"

	"github.com/foxcodenine/iot-parking-gateway/internal/core"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

// UDPServer represents a UDP server.
type UDPServer struct {
	Addr       string
	Connection *net.UDPConn
	app        *core.App
}

// NewServer initializes a new UDP server.
func NewServer(addr string, a *core.App) *UDPServer {
	return &UDPServer{Addr: addr, app: a}
}

// Start initializes and listens on the specified UDP address.
func (s *UDPServer) Start() error {
	udpAddr, err := net.ResolveUDPAddr("udp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	s.Connection, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("failed to start UDP server: %w", err)
	}
	helpers.LogInfo("UDP server started on %s\n", s.Addr)

	go s.listen() // Start listening in a goroutine

	return nil
}

// listen listens for incoming UDP messages.
func (s *UDPServer) listen() {
	defer s.Connection.Close()

	buffer := make([]byte, 1024)
	for {
		n, addr, err := s.Connection.ReadFromUDP(buffer)
		if err != nil {
			helpers.LogError(err, "Error reading UDP message")
			continue
		}
		go s.handleUDPMessage(s.Connection, buffer[:n], addr)
	}
}

// Stop gracefully stops the UDP server.
func (s *UDPServer) Stop() error {
	return s.Connection.Close()
}
