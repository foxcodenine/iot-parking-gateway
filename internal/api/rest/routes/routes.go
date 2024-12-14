package routes

import (
	"net/http"

	"path/filepath"
	"strings"

	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Routes sets up all HTTP routes and returns a router
func Routes() http.Handler {
	mux := chi.NewRouter()

	// Middleware
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	vueHandler := &handlers.VueHandler{}

	// Serve modified index.html for SPA routes
	mux.Get("/", vueHandler.ServeIndexWithVariables)
	mux.Get("/login", vueHandler.ServeIndexWithVariables)
	mux.Get("/forgot-password", vueHandler.ServeIndexWithVariables)
	mux.Get("/device", vueHandler.ServeIndexWithVariables)
	mux.Get("/user", vueHandler.ServeIndexWithVariables)
	mux.Get("/user/{id}", vueHandler.ServeIndexWithVariables)
	mux.Get("/app", vueHandler.ServeIndexWithVariables) // add more routes as needed

	// API routes
	mux.Route("/api", func(r chi.Router) {
		r.Mount("/device", DeviceRoutes())
		r.Mount("/user", UserRoutes())
		r.Mount("/auth", AuthRoutes())
		r.Mount("/setting", SettingRoutes())
	})

	// Serve all static files under the dist directory
	workDir, _ := filepath.Abs(".")
	filesDir := filepath.Join(workDir, "dist")
	FileServer(mux, "/", http.Dir(filesDir))

	return mux
}

// FileServer sets up a http.FileServer handler to serve static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	r.Get(path+"*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}
