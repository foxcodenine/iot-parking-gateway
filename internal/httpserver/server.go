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

type Server struct {
	HTTPServer *http.Server
}

func NewServer(port string) *Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: routes.Routes(), // Ensure this accesses the routes correctly
	}

	return &Server{
		HTTPServer: srv,
	}
}

func (s *Server) Start() {
	go func() {
		helpers.LogInfo("HTTP server starting on %s\n", s.HTTPServer.Addr)
		if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			helpers.LogError(err, "Error starting the server: %v\n")
		}
	}()
}

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
