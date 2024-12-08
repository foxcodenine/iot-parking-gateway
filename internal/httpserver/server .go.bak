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
)

// Server wraps the http.Server.
type Server struct {
	HTTPServer *http.Server
}

// NewServer initializes a new HTTP server on the specified port with routes configured.
func NewServer(port string) *Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: routes.Routes(), // Setup HTTP routing.
	}
	return &Server{
		HTTPServer: srv,
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
