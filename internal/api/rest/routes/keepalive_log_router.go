package routes

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/go-chi/chi/v5"
)

func KeepaliveLogRouter() chi.Router {
	r := chi.NewRouter()

	keepaliveLogHandler := &handlers.KeepaliveLogHandler{}

	// r.Use(middleware.JWTAuthMiddleware)

	r.Get("/{device_id}", keepaliveLogHandler.Get)

	return r
}
