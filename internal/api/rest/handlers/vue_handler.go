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

// Custom handler to inject variables into index.html
func (v *VueHandler) ServeIndexWithVariables(w http.ResponseWriter, r *http.Request) {
	// Define your dynamic variable
	data := struct {
		GoogleApiKey string
	}{
		GoogleApiKey: os.Getenv("GOOGLE_API_KEY"),
	}

	// Path to the index.html file
	filePath := filepath.Join("dist", "index.html")

	// Parse the HTML file
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		http.Error(w, "Error loading index.html", http.StatusInternalServerError)
		return
	}

	// Write the modified content to a buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Serve the modified content
	w.Header().Set("Cache-Control", "no-store, must-revalidate")
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
