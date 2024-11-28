package handlers

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type VueHandler struct {
}

func (v *VueHandler) ServeIndexWithVariables(w http.ResponseWriter, r *http.Request) {
	// Define your dynamic variable

	data := struct {
		GoogleApiKey string
	}{
		GoogleApiKey: os.Getenv("GOOGLE_API_KEY"),
	}

	// Path to the index.html file
	filePath := filepath.Join("dist", "index.html")

	// Read and parse the HTML file anew on each request
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		http.Error(w, "Error loading index.html", http.StatusInternalServerError)
		return
	}

	// Execute the template with the data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
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
