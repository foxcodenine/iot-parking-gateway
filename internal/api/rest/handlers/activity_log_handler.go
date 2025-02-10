package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/go-chi/chi/v5"
)

// ActivityLogHandler handles activity log requests.
type ActivityLogHandler struct{}

// Get handles the request for activity logs.
func (h *ActivityLogHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters
	deviceID := chi.URLParam(r, "device_id")
	deviceID = strings.TrimSpace(deviceID)
	fromDateStr := r.URL.Query().Get("from_date")
	toDateStr := r.URL.Query().Get("to_date")

	// Validate device ID
	if deviceID == "" {
		http.Error(w, "DeviceID cannot be empty", http.StatusBadRequest)
		return
	}

	// Convert timestamps from string to int64
	fromDate, err := strconv.ParseInt(fromDateStr, 10, 64)
	if err != nil || fromDate <= 0 {
		http.Error(w, "Invalid from_date. Must be a valid timestamp.", http.StatusBadRequest)
		return
	}

	toDate, err := strconv.ParseInt(toDateStr, 10, 64)
	if err != nil || toDate <= 0 {
		http.Error(w, "Invalid to_date. Must be a valid timestamp.", http.StatusBadRequest)
		return
	}

	if fromDate > toDate {
		http.Error(w, "from_date cannot be greater than to_date", http.StatusBadRequest)
		return
	}

	// Check if the device exists in cache
	deviceExists, err := cache.AppCache.Exists("parking:device:" + deviceID)
	if err != nil {
		helpers.RespondWithError(w, err, "Error checking device existence", http.StatusInternalServerError)
		return
	}
	if !deviceExists {
		helpers.RespondWithError(w, err, "Device does not exist", http.StatusNotFound)
		return
	}

	// Fetch activity logs from storage
	activityLogs, err := app.Models.ActivityLog.GetActivityLogs(deviceID, fromDate, toDate)
	if err != nil {
		helpers.RespondWithError(w, err, "Error fetching activity logs", http.StatusInternalServerError)
		return
	}

	// Determine the response message based on the number of logs retrieved
	var message string
	logCount := len(activityLogs)

	if logCount > 0 {
		message = fmt.Sprintf("%d activity logs retrieved successfully.", logCount)
	} else {
		message = "No activity logs found for the given asset within the specified date range."
	}

	// Response structure
	response := map[string]interface{}{
		"message":       message,
		"activity_logs": activityLogs,
	}

	// Encode and send the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
