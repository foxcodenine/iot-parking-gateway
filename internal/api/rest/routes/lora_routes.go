package routes

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"

	"github.com/go-chi/chi/v5"
)

func LoraRoutes() chi.Router {
	r := chi.NewRouter()

	loraHandler := &handlers.LoraHandler{}

	r.Post("/chirpstack", loraHandler.UpChirpstack)

	return r
}
