package routes

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/go-chi/chi/v5"
)

func DeviceRoutes() chi.Router {
	r := chi.NewRouter()

	deviceHandler := &handlers.DeviceHandler{}

	r.Get("/", deviceHandler.ListDevices)
	r.Get("/{id}", deviceHandler.GetDevice)
	r.Post("/", deviceHandler.CreateDevice)
	r.Put("/{id}", deviceHandler.UpdateDevice)
	r.Delete("/{id}", deviceHandler.DeleteDevice)

	return r
}
