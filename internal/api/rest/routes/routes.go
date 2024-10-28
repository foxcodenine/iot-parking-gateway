package routes

import (
	"net/http"

	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Routes sets up all HTTP routes and returns a router
func Routes() http.Handler {
	mux := chi.NewRouter()

	// Middleware
	mux.Use(middleware.Recoverer)

	// Retrieve the repository to access shared application configurations
	repo := handlers.GetRepo()

	// Initialize specific handlers using the repository
	testHandler := &handlers.TestHandler{Repo: repo}

	// Define routes for each handler
	mux.Get("/test", testHandler.Index)

	return mux
}
