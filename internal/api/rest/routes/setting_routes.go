package routes

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/middleware"
	"github.com/go-chi/chi/v5"
)

func SettingRoutes() chi.Router {
	r := chi.NewRouter()

	settingHandler := &handlers.SettingHandler{}

	r.Use(middleware.JWTAuthMiddleware)

	r.Put("/", settingHandler.Update)

	return r
}
