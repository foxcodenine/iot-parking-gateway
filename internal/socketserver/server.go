package socketserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/go-chi/chi/v5"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	gorillaWebSocket "github.com/gorilla/websocket"
)

type Server struct {
	SocketServer *socketio.Server
	SocketPort   string
	HttpServer   *http.Server
}

// allowOriginFunc allows all origins; used to configure CORS in the transports.
var allowOriginFunc = func(r *http.Request) bool {
	return true
}

var IOService *Server

// NewServer creates and configures a new Socket.IO server.
func NewServer(port string) *Server {
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
		helpers.LogInfo("Connected ID: %s", s.ID())
		return nil
	})

	socketServer.OnEvent("/", "update", func(s socketio.Conn, msg string) {
		helpers.LogInfo("Received update: %s", msg)
	})

	socketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		helpers.LogInfo("Disconnected ID: %s, Reason: %s", s.ID(), reason)
	})

	socketServer.OnError("/", func(s socketio.Conn, err error) {
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

	mux := chi.NewMux()
	mux.Handle("/socket.io/", socketServer)

	IOService = &Server{
		SocketServer: socketServer,
		SocketPort:   port,
		HttpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: mux,
		},
	}

	return IOService
}

// Start initiates the Socket.IO and HTTP servers in separate goroutines.
func (s *Server) Start() {
	helpers.LogInfo("Socket.IO server starting on port %s", s.SocketPort)

	go func() {
		if err := s.SocketServer.Serve(); err != nil {
			helpers.LogError(err, "Error starting Socket.IO server")
		}
	}()

	go func() {
		if err := s.HttpServer.ListenAndServe(); err != http.ErrServerClosed {
			helpers.LogError(err, "HTTP server stopped unexpectedly")
		}
	}()
}

// Shutdown handles graceful shutdown of the Socket.IO and HTTP servers.
func (s *Server) Shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	helpers.LogInfo("Shutting down Socket.IO server...")
	if err := s.SocketServer.Close(); err != nil {
		helpers.LogError(err, "Failed to shut down Socket.IO server")
	}

	helpers.LogInfo("Shutting down HTTP server...")
	if err := s.HttpServer.Shutdown(ctx); err != nil {
		helpers.LogError(err, "Failed to shut down HTTP server")
	}
}
