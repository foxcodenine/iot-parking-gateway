package routes

import (
	"os"

	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/go-chi/chi/v5"
)

func EnvRoutes() chi.Router {
	r := chi.NewRouter()

	// Initialize specific handlers using the repository
	envHandler := &handlers.EnvHandler{SecretKey: os.Getenv("SECRET_KEY")}

	r.Get("/", envHandler.Index)

	return r
}
