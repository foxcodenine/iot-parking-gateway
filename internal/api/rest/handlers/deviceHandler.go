package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

type DeviceHandler struct {
}

func (h *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) {

	devices, err := app.Models.Device.GetAll()

	if err != nil {
		app.ErrorLog.Printf("Failed to retrieve devices: %v", err)
		helpers.RespondWithError(w, err, "Failed to retrieve devices", http.StatusInternalServerError)
		return
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the devices slice to JSON and write it to the response
	err = json.NewEncoder(w).Encode(devices)

	if err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {

	var payload struct {
		IMEI string `json:"imei"`
	}

	// Decode the JSON body into the payload struct
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.ErrorLog.Printf("Failed to create device: %v", err)
		helpers.RespondWithError(w, err, "Failed to create device", http.StatusInternalServerError)
		return
	}

	// Call the Create method on the Device model (example uses hardcoded device ID)
	device, err := app.Models.Device.Create(payload.IMEI)
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

}

func (h *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {

}
func (h *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {

}
