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
	r.Get("/{id}", deviceHandler.GetDevice)
	r.Post("/", deviceHandler.CreateDevice)
	r.Put("/{id}", deviceHandler.UpdateDevice)
	r.Delete("/{id}", deviceHandler.DeleteDevice)

	return r
}
