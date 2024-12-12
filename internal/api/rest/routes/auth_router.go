package routes

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/middleware"
	"github.com/go-chi/chi/v5"
)

func AuthRoutes() chi.Router {
	r := chi.NewRouter()

	authHandler := &handlers.AuthHandler{}

	r.Post("/login", authHandler.Login)
	r.With(middleware.JWTAuthMiddleware).Post("/logout", authHandler.Logout)

	return r
}
