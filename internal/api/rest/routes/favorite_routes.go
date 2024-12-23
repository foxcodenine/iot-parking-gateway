package routes

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/middleware"
	"github.com/go-chi/chi/v5"
)

func FavoriteRoutes() chi.Router {
	r := chi.NewRouter()

	favoriteHandler := &handlers.FavoriteHandler{}

	r.Use(middleware.JWTAuthMiddleware)

	r.Put("/", favoriteHandler.Update)

	return r
}
