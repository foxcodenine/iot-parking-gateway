package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"

	"github.com/go-chi/chi/v5"
)

type DeviceHandler struct {
}

func (h *DeviceHandler) Index(w http.ResponseWriter, r *http.Request) {

	devices, err := app.Models.Device.GetAll()

	if err != nil {
		helpers.RespondWithError(w, err, "Failed to retrieve devices", http.StatusInternalServerError)
		return
	}

	// Response structure with a success message and user data
	response := map[string]interface{}{
		"message": "Devices retrieved successfully.",
		"devices": devices,
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the devices slice to JSON and write it to the response

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response.", http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandler) Store(w http.ResponseWriter, r *http.Request) {
	userData, err := app.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Authentication error.", http.StatusUnauthorized)
		return
	}

	// Check if the user has permission to update a user
	if userData.AccessLevel > 2 {
		http.Error(w, "You do not have the necessary permissions to perform this action.", http.StatusForbidden)
		return
	}

	var payload struct {
		DeviceID        string  `json:"device_id"`
		Name            string  `json:"name"`
		NetworkType     string  `json:"network_type"`
		FirmwareVersion float64 `json:"firmware_version"`
		Latitude        float64 `json:"latitude"`
		Longitude       float64 `json:"longitude"`
		IsAllowed       bool    `json:"is_allowed"`
		IsBlocked       bool    `json:"is_blocked"`
		IsHidden        bool    `json:"is_hidden"`
	}

	// Decode the JSON body into the payload struct
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to create device!!", http.StatusInternalServerError)
		return
	}

	// Validation logic
	if len(payload.DeviceID) < 3 || len(payload.Name) < 3 {
		http.Error(w, "Device ID and Name must be at least 3 characters long.", http.StatusBadRequest)
		return
	}

	validNetworkTypes := map[string]bool{"LoRa": true, "SigFox": true, "NB-IoT": true}
	if _, ok := validNetworkTypes[payload.NetworkType]; !ok {
		http.Error(w, "Network Type must be either 'LoRa', 'SigFox', or 'NB-IoT'.", http.StatusBadRequest)
		return
	}

	newDevice := models.Device{
		DeviceID:        payload.DeviceID,
		NetworkType:     payload.NetworkType,
		Name:            payload.Name,
		FirmwareVersion: payload.FirmwareVersion,
		Latitude:        payload.Latitude,
		Longitude:       payload.Longitude,
		IsAllowed:       payload.IsAllowed,
		IsBlocked:       payload.IsBlocked,
		IsHidden:        payload.IsHidden,
	}

	// Call the Create method on the Device model (example uses hardcoded device ID)
	device, err := app.Models.Device.Upsert(&newDevice)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to create device!!!", http.StatusInternalServerError)
		return
	}

	deviceIdentifierKey := fmt.Sprintf("%s %s", newDevice.NetworkType, newDevice.DeviceID)

	_, err = cache.AppCache.AddItemToBloomFilter("registered-devices", deviceIdentifierKey)
	if err != nil {
		helpers.LogError(err, "Failed to add device ID to the Bloom Filter")
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	app.PushAuditToCache(*userData, "UPDATE", "device", newDevice.DeviceID, r, fmt.Sprintf("Created device with ID %s.", newDevice.DeviceID))

	// Encode the created device to JSON and write it to the response
	if err = json.NewEncoder(w).Encode(device); err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Attempt to retrieve the device by ID
	device, err := app.Models.Device.GetByID(id)

	if err != nil {
		app.ErrorLog.Printf("Failed to retrieve device %s: %v", id, err)
		helpers.RespondWithError(w, err, "Failed to retrieve device", http.StatusInternalServerError)
		return
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the created device to JSON and write it to the response
	if err = json.NewEncoder(w).Encode(device); err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandler) Update(w http.ResponseWriter, r *http.Request) {
	userData, err := app.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Authentication error.", http.StatusUnauthorized)
		return
	}

	// Check if the user has permission to update a user
	if userData.AccessLevel > 2 {
		http.Error(w, "You do not have the necessary permissions to perform this action.", http.StatusForbidden)
		return
	}

	// Get the device ID from the URL
	id := chi.URLParam(r, "id")

	// Define a map to hold the fields to update
	var updatedFields map[string]interface{}

	// Decode the JSON body into the map for flexibility
	err = json.NewDecoder(r.Body).Decode(&updatedFields)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Attempt to update the device
	device, err := app.Models.Device.UpdateByID(id, updatedFields)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to update device", http.StatusInternalServerError)
		return
	}

	app.PushAuditToCache(*userData, "UPDATE", "device", id, r, fmt.Sprintf("Updated device with ID %s.", id))

	// Response structure with a success message and user data
	response := map[string]interface{}{
		"message": "Devices updated successfully.",
		"device":  device,
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the updated device to JSON and write it to the response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	userData, err := app.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Authentication error.", http.StatusUnauthorized)
		return
	}

	// Check if the user has sufficient permission to delete a device
	if userData.AccessLevel > 2 {
		http.Error(w, "You do not have the necessary permissions to perform this action.", http.StatusForbidden)
		return
	}

	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "Device ID is required.", http.StatusBadRequest)
		return
	}

	// Attempt to soft delete the device
	err = app.Models.Device.SoftDeleteByID(id)
	if err != nil {
		app.ErrorLog.Printf("Failed to soft delete device %s: %v", id, err)
		if err.Error() == "device not found" {
			http.Error(w, fmt.Sprintf("Device with ID %s not found.", id), http.StatusNotFound)
		} else {
			helpers.RespondWithError(w, err, "Failed to delete device.", http.StatusInternalServerError)
		}
		return
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	app.PushAuditToCache(*userData, "DELETE", "device", id, r, fmt.Sprintf("Marked device with ID %s as soft deleted.", id))

	response := map[string]string{"message": "Device deleted successfully"}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
