package routes

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/middleware"
	"github.com/go-chi/chi/v5"
)

func UserRoutes() chi.Router {
	r := chi.NewRouter()

	userHandler := &handlers.UserHandler{}

	r.Use(middleware.JWTAuthMiddleware)

	r.Post("/", userHandler.CreateUser)

	return r
}
