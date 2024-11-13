package udp

import (
	"net"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

// UDPServer represents a UDP server.
type UDPServer struct {
	Addr           string
	Connection     *net.UDPConn
	cache          *cache.RedisCache
	shutdownCh     chan struct{} // Shutdown channel to signal the listening loop to stop
	isShuttingDown bool          // Flag to indicate the server is intentionally shutting down

}

// NewServer initializes a new UDP server.
func NewServer(addr string, c *cache.RedisCache) *UDPServer {
	return &UDPServer{
		Addr:           addr,
		cache:          c,
		shutdownCh:     make(chan struct{}),
		isShuttingDown: false,
	}
}

// Start initializes and listens on the specified UDP address.
func (s *UDPServer) Start() {
	udpAddr, err := net.ResolveUDPAddr("udp", s.Addr)
	if err != nil {
		helpers.LogFatal(err, "Failed to resolve UDP address")
	}

	s.Connection, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		helpers.LogFatal(err, "Failed to start UDP server")
	}
	helpers.LogInfo("UDP server started on %s\n", s.Addr)

	go s.listen() // Start listening in a goroutine

}

// listen listens for incoming UDP messages.
func (s *UDPServer) listen() {
	defer s.Connection.Close()
	buffer := make([]byte, 1024)

	for {
		select {
		case <-s.shutdownCh: // Listen for shutdown signal.
			helpers.LogInfo("Shutdown signal received, stopping UDP listener")
			return // Properly handle the shutdown signal by exiting the loop.
		default:
			n, addr, err := s.Connection.ReadFromUDP(buffer)
			if err != nil {
				if err == net.ErrClosed {
					helpers.LogError(err, "UDP connection unexpectedly closed!")
					return // Exit the loop gracefully
				}
				if !s.isShuttingDown {
					helpers.LogError(err, "Error reading UDP message")
				}
				continue // Handle other errors and continue listening
			}

			// Process the received data in a goroutine
			go s.nbMessageHandler(s.Connection, buffer[:n], addr)
		}
	}
}

// Stop gracefully stops the UDP server.
func (s *UDPServer) Stop() {
	helpers.LogInfo("Initiating shutdown of the UDP server...")
	s.isShuttingDown = true
	close(s.shutdownCh)         // Signal to stop the listening loop
	err := s.Connection.Close() // Then close the connection

	if err != nil {
		helpers.LogError(err, "Failed to gracefully stop UDP server")
	}
	helpers.LogInfo("UDP connection closed.")

}
