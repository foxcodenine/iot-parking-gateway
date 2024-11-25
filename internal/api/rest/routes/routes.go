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

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Initialize specific handlers using the repository
	testHandler := &handlers.TestHandler{}
	vueHandler := &handlers.VueHandler{}

	// Define routes for each handler
	mux.Get("/test", testHandler.Index)

	// -----------------------------------------------------------------

	// Serve static index.html at route /app
	// mux.Get("/", vueHandler.ServeIndexWithVariables)
	// mux.Get("/login", vueHandler.ServeIndexWithVariables)
	// mux.Get("/user", vueHandler.ServeIndexWithVariables)

	// Handle all other routes

	// -----------------------------------------------------------------

	// Mount api routes
	mux.Route("/api", func(r chi.Router) {
		r.Mount("/device", DeviceRoutes())
		r.Mount("/user", UserRoutes())
		r.Mount("/auth", AuthRoutes())
	})

	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {

		// Serve all static files under the dist directory
		workDir, _ := filepath.Abs(".")
		filesDir := filepath.Join(workDir, "dist")
		FileServer(mux, "/", http.Dir(filesDir))

		if !strings.HasPrefix(r.URL.Path, "/api") {
			// Let the API handler handle it
			vueHandler.ServeIndexWithVariables(w, r)
		}
		http.NotFound(w, r)
		return
	})

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
