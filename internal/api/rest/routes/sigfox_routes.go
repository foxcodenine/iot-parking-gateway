package routes

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"

	"github.com/go-chi/chi/v5"
)

func SigfoxRoutes() chi.Router {
	r := chi.NewRouter()

	sigfoxHandler := &handlers.SigfoxHandler{}

	r.Post("/", sigfoxHandler.Up)

	return r
}
