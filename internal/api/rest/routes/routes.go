package routes

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Routes sets up all HTTP routes and returns a router
func Routes() http.Handler {
	mux := chi.NewRouter()

	// Middleware
	mux.Use(middleware.Recoverer)

	// Initialize specific handlers using the repository
	testHandler := &handlers.TestHandler{}

	// Define routes for each handler
	mux.Get("/test", testHandler.Index)

	// Serve static index.html at route /app
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join("dist", "index.html")
		http.ServeFile(w, r, filePath)
	})

	// Serve all static files under the dist directory
	workDir, _ := filepath.Abs(".")
	filesDir := filepath.Join(workDir, "dist")
	FileServer(mux, "/", http.Dir(filesDir))

	// Mount device routes
	mux.Route("/api", func(r chi.Router) {
		r.Mount("/device", DeviceRoutes())
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
