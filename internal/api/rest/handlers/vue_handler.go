package handlers

import (
	"bytes"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

type VueHandler struct {
}

func (v *VueHandler) ServeIndexWithVariables(w http.ResponseWriter, r *http.Request) {
	// Define your dynamic variable
	// googleApiKey, _ := helpers.EncryptAES(os.Getenv("GOOGLE_API_KEY"), core.AES_SECRET_KEY)

	loginPageTitle, err := cache.AppCache.HGet("app:settings", "login_page_title")
	if err != nil {
		// Log the error with a specific message to indicate fetching login page title failed
		helpers.LogError(err, "Failed to fetch login page title from cache; using default value.")
		// Set a default value if there's an error fetching the cached value
		loginPageTitle = "Welcome to <b>IoTrack</b> Pro"
	}

	// Safe type assertion
	loginPageTitleStr, ok := loginPageTitle.(string)
	if !ok {
		// Log or handle the unexpected type scenario
		helpers.LogError(nil, "Expected string for login page title but got another type; using default value.")
		loginPageTitleStr = "Welcome to <b>IoTrack</b> Pro"
	}

	data := struct {
		// GoogleApiKey string
		LoginPageTitle string
	}{
		LoginPageTitle: loginPageTitleStr,
	}

	// Path to the index.html file
	filePath := filepath.Join("dist", "index.html")

	// Read and parse the HTML file anew on each request
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		helpers.RespondWithError(w, err, "Error loading index.html", http.StatusInternalServerError)
		return
	}

	// Execute the template with the data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		helpers.RespondWithError(w, err, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Set headers to control caching and content type
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // Ensures no caching
	w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0 backwards compatibility
	w.Header().Set("Expires", "0")                                         // Proxies
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
