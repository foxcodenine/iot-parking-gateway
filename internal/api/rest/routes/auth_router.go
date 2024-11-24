package routes

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/go-chi/chi/v5"
)

func AuthRoutes() chi.Router {
	r := chi.NewRouter()

	authHandler := &handlers.AuthHandler{}

	r.Post("/login", authHandler.Login)

	return r
}
