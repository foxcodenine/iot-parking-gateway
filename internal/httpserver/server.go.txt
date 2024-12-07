// Package httpserver handles the setup and lifecycle of the HTTP and Socket.IO servers.

package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/routes"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

// Server encapsulates both the HTTP server and the Socket.IO server,
// managing their lifecycle and integration.
type Server struct {
	HTTPServer   *http.Server
	SocketServer *socketio.Server
}

// Easier to get running with CORS. Thanks for help @Vindexus and @erkie
var allowOriginFunc = func(r *http.Request) bool {
	return true
}

// NewServer initializes and returns a new Server instance, setting up
// both the HTTP and Socket.IO servers with their respective routes and handlers.
func NewServer(port string) *Server {
	// Initialize the Socket.IO server
	socketServer := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	socketServer.OnConnect("/", func(s socketio.Conn) error {
		fmt.Printf("New connection established: %s\n", s.ID()) // Add log here
		helpers.LogInfo("Connected ID: %s", s.ID())
		return nil
	})

	socketServer.OnEvent("/", "update", func(s socketio.Conn, msg string) {
		fmt.Printf("Event received from client: %s\n", msg) // Add log here
		helpers.LogInfo("Received update: %s", msg)
	})

	socketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Printf("Client disconnected: %s, Reason: %s\n", s.ID(), reason) // Add log here
		helpers.LogInfo("Disconnected ID: %s, Reason: %s", s.ID(), reason)
	})

	socketServer.OnError("/", func(s socketio.Conn, err error) {
		// Log any errors that occur within the Socket.IO server
		helpers.LogError(err, "sockert-io error")
	})

	// Set up the HTTP multiplexer to handle both HTTP routes and Socket.IO routes
	mux := http.NewServeMux()
	mux.Handle("/", routes.Routes())
	mux.Handle("/socket.io/", socketServer)

	// Create the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	// Return the Server instance
	return &Server{
		HTTPServer:   srv,
		SocketServer: socketServer,
	}
}

// Start begins the HTTP server in a separate goroutine.
func (s *Server) Start() {
	go func() {
		helpers.LogInfo("HTTP server starting on %s\n", s.HTTPServer.Addr)
		err := s.HTTPServer.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			helpers.LogError(err, "Error starting the server: %v\n")
		}
	}()
}

// Shutdown gracefully stops the HTTP and Socket.IO servers.
// It waits for any ongoing requests or connections to complete before shutting down.
func (s *Server) Shutdown() {

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until a signal is received

	// Create a context with a timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Log and initiate the HTTP server shutdown
	helpers.LogInfo("HTTP server shutting down...")
	err := s.HTTPServer.Shutdown(ctx)
	if err != nil {
		helpers.LogError(err, "Server forced to shutdown\n")
	}

	// Close the Socket.IO server
	err = s.SocketServer.Close()
	if err != nil {
		helpers.LogError(err, "Error closing Socket.IO server \n")
	}
	helpers.LogInfo("HTTP server shutdown complete.")
}
