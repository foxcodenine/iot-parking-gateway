package routes

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/middleware"
	"github.com/go-chi/chi/v5"
)

func ActivityLogRouter() chi.Router {
	r := chi.NewRouter()

	activityLogHandler := &handlers.ActivityLogHandler{}

	r.Use(middleware.JWTAuthMiddleware)

	// Use GET for retrieving activity logs with query parameters
	// eg: GET /activity-logs/02DF9902?from_date=1707803200&to_date=1707806000
	r.Get("/{device_id}", activityLogHandler.Get)

	return r
}
