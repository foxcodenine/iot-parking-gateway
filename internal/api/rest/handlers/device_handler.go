package handlers

import (
	"encoding/json"
	"net/http"

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

func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {

	var payload struct {
		DeviceID    string `json:"device_id"`
		NetworkType string `json:"network_type"`
	}

	// Decode the JSON body into the payload struct
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.ErrorLog.Printf("Failed to create device: %v", err)
		helpers.RespondWithError(w, err, "Failed to create device", http.StatusInternalServerError)
		return
	}

	newDevice := models.Device{
		DeviceID:    payload.DeviceID,
		NetworkType: payload.NetworkType,
	}

	// Call the Create method on the Device model (example uses hardcoded device ID)
	device, err := app.Models.Device.Create(&newDevice)
	if err != nil {
		app.ErrorLog.Printf("Failed to create device: %v", err)
		helpers.RespondWithError(w, err, "Failed to create device", http.StatusInternalServerError)
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

func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
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

func (h *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	// Get the device ID from the URL
	id := chi.URLParam(r, "id")

	// Define a map to hold the fields to update
	var updatedFields map[string]interface{}

	// Decode the JSON body into the map for flexibility
	err := json.NewDecoder(r.Body).Decode(&updatedFields)
	if err != nil {
		app.ErrorLog.Printf("Failed to decode request body: %v", err)
		helpers.RespondWithError(w, err, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Attempt to update the device
	device, err := app.Models.Device.UpdateByID(id, updatedFields)
	if err != nil {
		app.ErrorLog.Printf("Failed to update device %s: %v", id, err)
		helpers.RespondWithError(w, err, "Failed to update device", http.StatusInternalServerError)
		return
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the updated device to JSON and write it to the response
	err = json.NewEncoder(w).Encode(device)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := app.Models.Device.DeleteByID(id)

	if err != nil {
		app.ErrorLog.Printf("Failed to delete device %s: %v", id, err)
		helpers.RespondWithError(w, err, "Failed to delete device", http.StatusInternalServerError)
		return
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{"message": "Device deleted successfully"}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
