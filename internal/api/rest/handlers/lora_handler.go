package handlers

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	lorafw "github.com/foxcodenine/iot-parking-gateway/internal/firmware/lora_fw"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/foxcodenine/iot-parking-gateway/internal/mq"

	"github.com/google/uuid"
)

type LoraHandler struct {
}

func (h *LoraHandler) UpChirpstack(w http.ResponseWriter, r *http.Request) {
	type DeviceInfo struct {
		TenantId          string `json:"tenantId"`
		TenantName        string `json:"tenantName"`
		ApplicationId     string `json:"applicationId"`
		ApplicationName   string `json:"applicationName"`
		DeviceProfileId   string `json:"deviceProfileId"`
		DeviceProfileName string `json:"deviceProfileName"`
		DeviceName        string `json:"deviceName"`
		DevEui            string `json:"devEui"`
	}

	type Request struct {
		DeviceInfo DeviceInfo `json:"deviceInfo"`
		Data       string     `json:"data"`
	}
	var req Request

	// Parse the JSON payload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.RespondWithError(w, err, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// -----------------------------------------------------------------

	// Base64 encoded string
	base64Str := req.Data

	// Decode the base64 string
	bufferBase64, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		helpers.RespondWithError(w, err, "Error decoding base64", http.StatusInternalServerError)
		return
	}

	// Convert the decoded bytes to a hex string
	hexStr := hex.EncodeToString(bufferBase64)

	// -----------------------------------------------------------------

	// hexStr := req.Data

	// -----------------------------------------------------------------

	// Parse firmware version
	firmwareVersionTmp, _, err := helpers.ParseHexSubstring(hexStr, 0, 1)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to parse firmware version", http.StatusInternalServerError)
		return
	}

	// Divide by 10 to convert to float64 and shift decimal place
	firmwareVersion := float64(firmwareVersionTmp) / 10.0
	deviceID := strings.ToUpper(req.DeviceInfo.DeviceName)
	rawData := req
	rawData.Data = hexStr

	// Convert struct to JSON string
	rawDataBytes, err := json.Marshal(rawData)
	if err != nil {
		helpers.RespondWithError(w, err, "Error converting to JSON String", http.StatusInternalServerError)
		return

	}
	rawDataString := string(rawDataBytes)

	// Check if the device ID is already in the Bloom Filter
	deviceIdentifierKey := fmt.Sprintf("LoRa %s", deviceID)
	isDeviceRegistered, err := cache.AppCache.CheckItemInBloomFilter("registered-devices", deviceIdentifierKey)
	if err != nil {
		helpers.LogError(err, "Failed to check Bloom Filter for device ID")
	}

	// If the device ID is not registered, track it for registration and prevent duplicates.
	if !isDeviceRegistered {
		// Add the device to a Redis set for tracking devices that need registration.
		deviceDataKey := fmt.Sprintf("%s %f", deviceIdentifierKey, firmwareVersion)
		if err := cache.AppCache.SAdd("to-register-devices", deviceDataKey); err != nil {
			helpers.LogError(err, "Failed to add device ID to the 'to-register-devices' set")
		}

		// Add the device ID to the Bloom Filter to prevent duplicate registrations in the future.
		if _, err := cache.AppCache.AddItemToBloomFilter("registered-devices", deviceIdentifierKey); err != nil {
			helpers.LogError(err, "Failed to add device ID to the 'registered-devices' Bloom Filter")
		}
	}

	// If the device ID is registered, check if it is soft delete, white listed or black listed.
	if isDeviceRegistered {

		// Retrieve the device data from the cache
		deviceData, err := cache.AppCache.GetDevice(deviceID)
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to retrieve device data from cache.", http.StatusInternalServerError)
			return
		}

		// Check if the device is soft deleted
		if deletedAt, exists := deviceData["deleted_at"]; exists && deletedAt != nil && deletedAt != "0001-01-01T00:00:00Z" {
			response := map[string]interface{}{
				"status":  "ignored",
				"message": fmt.Sprintf("Device %s is marked as soft deleted. Request ignored.", deviceID),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusSeeOther)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Retrieve application settings for device access mode
		deviceAccessMode, err := cache.AppCache.HGet("app:settings", "device_access_mode")
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to retrieve 'device_access_mode' from application settings.", http.StatusInternalServerError)
			return
		}

		// Check access based on whitelist mode
		if deviceAccessMode == "white_list" {
			if isAllowed, ok := deviceData["is_allowed"].(bool); ok && !isAllowed {
				response := map[string]interface{}{
					"status":  "ignored",
					"message": fmt.Sprintf("Device %s is not marked allowed. Request ignored.", deviceID),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden) // 403 Forbidden
				json.NewEncoder(w).Encode(response)
				return
			}
		}

		// Check access based on blacklist mode
		if deviceAccessMode == "black_list" {
			if isBlocked, ok := deviceData["is_blocked"].(bool); ok && isBlocked {
				response := map[string]interface{}{
					"status":  "ignored",
					"message": fmt.Sprintf("Device %s is marked blocked. Request ignored.", deviceID),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden) // 403 Forbidden
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}

	// Generate a new UUID for the RawDataLog entry
	rawUUID, err := uuid.NewV7()
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to generate UUID for RawDataLog entry.", http.StatusInternalServerError)
		return
	}

	// Create a new RawDataLog object to store in Redis.
	rawDataLog := models.RawDataLog{
		ID:              rawUUID,
		DeviceID:        deviceID,
		FirmwareVersion: firmwareVersion,
		NetworkType:     "LoRa",
		RawData:         rawDataString,
		CreatedAt:       time.Now(),
	}

	// Push the raw data log entry to Redis
	err = cache.AppCache.RPush("logs:raw-data-logs", rawDataLog)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to push raw data log to Redis.", http.StatusInternalServerError)
		return
	}

	// Debug output for parsed values
	helpers.LogInfo("Network: LoRa, Firmware: %.2f, Device ID: %s", firmwareVersion, deviceID)

	// Process firmware-specific data parsing based on the firmware version.
	var parsedData map[string]any

	switch firmwareVersion {
	case 5.8, 5.9:
		parsedData, err = lorafw.Lora_58(hexStr)

	default:

		response := map[string]interface{}{
			"status":  "unsupported_firmware",
			"message": fmt.Sprintf("Device %s has an unsupported firmware version:  %.2f. Request ignored.", deviceID, firmwareVersion),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusSeeOther)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err != nil {
		helpers.RespondWithError(w, err, fmt.Sprintf("Failed to parse data from Lora_%.0f firmware", firmwareVersion*10), http.StatusInternalServerError)
		return
	}

	helpers.PrettyPrintJSON(parsedData)

	// Push parsed parking data packages to Redis.
	for _, i := range parsedData["parking_packages"].([]map[string]any) {

		i["firmware_version"] = parsedData["firmware_version"]
		i["device_id"] = deviceID
		i["raw_id"] = rawUUID
		i["event_id"] = 26
		i["network_type"] = "LoRa"

		err := cache.AppCache.RPush("logs:activity-logs", i)
		if err != nil {
			helpers.LogError(err, "Failed to push parking package data log to Redis")
		}

		messageData, err := json.Marshal(i)
		if err != nil {
			helpers.LogError(err, "Failed to serialize parsedData to JSON")
			continue
		}

		mq.AppRabbitMQProducer.SendMessage("event_logs_exchange", "event_logs_queue", string(messageData))
	}
}
