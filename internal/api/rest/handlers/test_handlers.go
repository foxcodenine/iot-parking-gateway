package handlers

import (
	"fmt"
	"net/http"
)

// TestHandler is a handler for test routes, with access to shared Repository
type TestHandler struct {
}

// Index handles the /test route
func (h *TestHandler) Index(w http.ResponseWriter, r *http.Request) {
	repo.App.InfoLog.Println("Test Route") // Log to console for server-side tracking
	fmt.Fprint(w, "Test Route")            // Write response to client
}
