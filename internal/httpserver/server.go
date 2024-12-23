// internal/httpserver/server.go
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

	// "github.com/go-chi/chi/v5"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	gorillaWebSocket "github.com/gorilla/websocket"
)

// Server wraps the http.Server.
type Server struct {
	HTTPServer   *http.Server
	SocketServer *socketio.Server
	Port         string
}

var allowedOrigins = map[string]bool{
	"https://yourdomain.com":    true,
	"https://anotherdomain.com": true,
	"http://localhost:5173":     true,
	"http://127.0.0.1:5173":     true,
}

// allowOriginFunc allows all origins; used to configure CORS in the transports.
var allowOriginFunc = func(r *http.Request) bool {
	// return true
	return allowedOrigins[r.Header.Get("Origin")]
}

var IOService *Server

// NewServer initializes a new HTTP server on the specified port with routes configured.
func NewServer(port string) *Server {

	socketServer := createSocketServer()

	mux := http.NewServeMux()
	mux.Handle("/socket.io/", socketServer)
	mux.Handle("/", routes.Routes())

	// ------------------------------------------
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	return &Server{
		HTTPServer:   srv,
		SocketServer: socketServer,
		Port:         port,
	}
}

// Start begins running the HTTP server in a separate goroutine to allow it to listen for incoming requests without blocking the main thread.
func (s *Server) Start() {
	go func() {
		helpers.LogInfo("HTTP server starting on %s\n", s.HTTPServer.Addr)
		if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			helpers.LogError(err, "Error starting the server: %v\n")
		}
	}()

	s.registerSocketHandlers()

	go func() {
		if err := s.SocketServer.Serve(); err != nil {
			helpers.LogError(err, "Error starting Socket.IO server")
		}
	}()

}

func (s *Server) RestartSocketServer() {
	helpers.LogInfo("Restarting Socket.IO server...")

	// Close the existing Socket.IO server to disconnect all clients.
	if err := s.SocketServer.Close(); err != nil {
		helpers.LogError(err, "Error closing Socket.IO server")
	}

	// Reinitialize the Socket.IO server
	s.SocketServer = createSocketServer()

	s.registerSocketHandlers()

	// Restart Socket.IO server serve
	go func() {
		if err := s.SocketServer.Serve(); err != nil {
			helpers.LogError(err, "Error starting Socket.IO server after restart")
		}
	}()
	helpers.LogInfo("Socket.IO server restarted successfully.")
}

func (s *Server) registerSocketHandlers() {

	// Gracefully handles panics during socket handler registration,
	// logging the error without crashing the server.
	defer func() {
		if r := recover(); r != nil {
			helpers.LogError(fmt.Errorf("panic in socket handler registration: %v", r), "Handler registration failed")
		}
	}()

	s.SocketServer.OnConnect("/", func(s socketio.Conn) error {
		helpers.LogInfo("Connected ID: %s", s.ID())
		return nil
	})

	s.SocketServer.OnEvent("/", "update", func(s socketio.Conn, msg string) {
		helpers.LogInfo("Received update: %s", msg)
	})

	s.SocketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		helpers.LogInfo("Disconnected ID: %s, Reason: %s", s.ID(), reason)
	})

	s.SocketServer.OnError("/", func(s socketio.Conn, err error) {
		// Handle WebSocket-specific errors using gorillaWebSocket package
		if websocketErr, ok := err.(*gorillaWebSocket.CloseError); ok {
			switch websocketErr.Code {
			case gorillaWebSocket.CloseGoingAway, gorillaWebSocket.CloseNormalClosure:
				helpers.LogInfo("Normal disconnect by client ID: %s, Code: %d", s.ID(), websocketErr.Code)
				return
			}
		}
		helpers.LogError(err, "Socket.IO error")
	})
}

func createSocketServer() *socketio.Server {
	return socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
// It listens for SIGINT and SIGTERM signals to trigger a shutdown, allowing for graceful termination.
func (s *Server) Shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	helpers.LogInfo("HTTP server shutting down...")
	if err := s.HTTPServer.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}
}
