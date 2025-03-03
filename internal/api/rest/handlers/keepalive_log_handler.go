package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/go-chi/chi/v5"
)

// KeepaliveLogHandler handles HTTP requests related to keepalive logs.
type KeepaliveLogHandler struct{}

// Get retrieves keepalive logs for a given device within a specified date range.
func (h *KeepaliveLogHandler) Get(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "device_id")
	deviceID = strings.TrimSpace(deviceID)

	fromDateStr := r.URL.Query().Get("from_date")
	toDateStr := r.URL.Query().Get("to_date")

	if deviceID == "" || deviceID == "undefined" {
		http.Error(w, "DeviceID cannot be empty", http.StatusBadRequest)
		return
	}

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

	deviceExists, err := cache.AppCache.Exists("parking:device:" + deviceID)
	if err != nil {
		helpers.RespondWithError(w, err, "Error checking device existence", http.StatusInternalServerError)
		return
	}
	if !deviceExists {
		helpers.RespondWithError(w, err, "Device does not exist", http.StatusNotFound)
		return
	}

	cachedDevice, err := cache.AppCache.GetDevice(deviceID)
	if err != nil {
		helpers.RespondWithError(w, err, "Error retrieving device from cache", http.StatusInternalServerError)
		return
	}

	var message string
	var logs interface{}

	// Retrieve logs based on network type.
	switch cachedDevice["network_type"].(string) {
	case "LoRa":
		logs, err = (&models.LoraKeepaliveLog{}).GetActivityLogs(deviceID, fromDate, toDate)
	case "NB-IoT":
		logs, err = (&models.NbiotKeepaliveLog{}).GetActivityLogs(deviceID, fromDate, toDate)
	case "SigFox":
		logs, err = (&models.SigfoxKeepaliveLog{}).GetActivityLogs(deviceID, fromDate, toDate)
	default:
		helpers.RespondWithError(w, nil, "Unsupported network type", http.StatusBadRequest)
		return
	}

	if err != nil {
		helpers.RespondWithError(w, err, "Error retrieving keepalive logs", http.StatusInternalServerError)
		return
	}

	logCount := reflect.ValueOf(logs).Len() // Use reflection to count elements in interface{} that contains a slice
	if logCount > 0 {
		message = fmt.Sprintf("%d keepalive logs retrieved successfully.", logCount)
	} else {
		message = "No keepalive logs found for the given asset within the specified date range."
	}

	response := map[string]interface{}{
		"message":        message,
		"keepalive_logs": logs,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
