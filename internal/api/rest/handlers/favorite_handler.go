package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

type FavoriteHandler struct {
}

func (h *FavoriteHandler) Update(w http.ResponseWriter, r *http.Request) {

	userData, err := app.GetUserFromContext(r.Context())
	if err != nil {
		app.ErrorLog.Printf("Authentication error: %v", err)
		http.Error(w, "Authentication error.", http.StatusUnauthorized)
		return
	}

	// Check if the user has permission to update favorites
	const permissionLevelRequired = 3
	if userData.AccessLevel > permissionLevelRequired {
		http.Error(w, "You do not have the necessary permissions to perform this action.", http.StatusForbidden)
		return
	}

	// Parse and validate input from the API
	type Request struct {
		DeviceIDs []string `json:"device_ids"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}

	err = cache.AppCache.HSet("app:user:favorites", userData.ID, req.DeviceIDs)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to update favorites in Redis.", http.StatusInternalServerError)
		return
	}

	// Respond with success
	response := map[string]interface{}{
		"message":    "Favorites updated successfully.",
		"device_ids": req.DeviceIDs,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response.", http.StatusInternalServerError)
	}

}
