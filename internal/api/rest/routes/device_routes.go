package routes

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/middleware"
	"github.com/go-chi/chi/v5"
)

func DeviceRoutes() chi.Router {
	r := chi.NewRouter()

	deviceHandler := &handlers.DeviceHandler{}

	r.Use(middleware.JWTAuthMiddleware)

	r.Get("/", deviceHandler.Index)
	r.Get("/{id}", deviceHandler.Get)
	r.Post("/", deviceHandler.Store)
	r.Put("/{id}", deviceHandler.Update)
	r.Delete("/{id}", deviceHandler.Destroy)

	return r
}
