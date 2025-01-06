package handlers

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	sigfoxfw "github.com/foxcodenine/iot-parking-gateway/internal/firmware/sigfox_fw"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/foxcodenine/iot-parking-gateway/internal/mq"

	"github.com/google/uuid"
)

type SigfoxHandler struct {
}

func (h *SigfoxHandler) Up(w http.ResponseWriter, r *http.Request) {
	// Parse and validate input from the API
	type Request struct {
		Timestamp int    `json:"timestamp"`
		DeviceID  string `json:"device"`
		SeqNumber string `json:"seq_number"`
		Data      string `json:"data"`
	}

	var req Request

	// Parse the JSON payload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	// Validate minimum hex string length
	if len(hexStr) < 5 {
		http.Error(w, "invalid message length, incoming data too short for parsing", http.StatusBadRequest)
		return
	}

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

	deviceID := strings.ToUpper(req.DeviceID)
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
	deviceIdentifierKey := fmt.Sprintf("SigFox %s", deviceID)
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
		NetworkType:     "SigFox",
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
	helpers.LogInfo("Network: SigFox, Firmware: %.2f, Device ID: %s", firmwareVersion, deviceID)

	// Process firmware-specific data parsing based on the firmware version.
	var parsedData map[string]any
	switch firmwareVersion {
	case 6:
		parsedData, err = sigfoxfw.Sigfox_60(hexStr, req.Timestamp)

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
		helpers.RespondWithError(w, err, fmt.Sprintf("Failed to parse data from Sigfox_%.0f firmware", firmwareVersion*10), http.StatusInternalServerError)
		return
	}

	err = h.updateDeviceCacheAndBroadcast(parsedData, deviceID)

	// Check for errors in the update process.
	if err != nil {
		// Log the error with additional context for better troubleshooting.
		helpers.LogError(err, "Failed to update device cache and broadcast changes")
	}

	// Attempt to update device keepalive_at in cache and broadcast the changes.
	err = h.updateDeviceKeepaliveInCacheAndBroadcast(parsedData, deviceID)

	// Check for errors in the update process.
	if err != nil {
		// Log the error with additional context for better troubleshooting.
		helpers.LogError(err, "Failed to update device keepalive_at in cache and broadcast changes")
	}

	// Attempt to update device settings_at in cache, check if device_settings should be updated and broadcast settings_at.
	updateDeviceSettings, err := h.updateDeviceSettingsInCacheAndBroadcast(parsedData, deviceID)

	// Check for errors in the update process.
	if err != nil {
		// Log the error with additional context for better troubleshooting.
		helpers.LogError(err, "Failed to update device settings_at in cache and broadcast it")
	}

	helpers.PrettyPrintJSON(parsedData)

	// Push parsed parking data packages to Redis.
	for _, i := range parsedData["parking_packages"].([]map[string]any) {

		i["firmware_version"] = parsedData["firmware_version"]
		i["device_id"] = deviceID
		i["raw_id"] = rawUUID
		i["event_id"] = 26
		i["network_type"] = "SigFox"

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

	// Push parsed keepalive data to Redis.
	for _, i := range parsedData["keep_alive_packages"].([]map[string]any) {
		i["firmware_version"] = parsedData["firmware_version"]
		i["device_id"] = deviceID
		i["raw_id"] = rawUUID
		i["event_id"] = 6
		i["network_type"] = "SigFox"

		err := cache.AppCache.RPush("logs:sigfox-keepalive-logs", i)
		if err != nil {
			helpers.LogError(err, "Failed to push keepalive package data log to Redis")
		}

		messageData, err := json.Marshal(i)
		if err != nil {
			helpers.LogError(err, "Failed to serialize parsedData to JSON")
			continue
		}
		mq.AppRabbitMQProducer.SendMessage("event_logs_exchange", "event_logs_queue", string(messageData))
	}

	// Push parsed settings data to Redis.
	for n, i := range parsedData["settings_packages"].([]map[string]any) {
		if n == 0 && updateDeviceSettings {
			i["update_device_settings"] = true
		} else {
			i["update_device_settings"] = false
		}

		// Add common fields to each individual package
		i["firmware_version"] = parsedData["firmware_version"]
		i["device_id"] = deviceID
		i["raw_id"] = rawUUID
		i["event_id"] = 25 // Assuming 25 is the event ID for setting logs
		i["network_type"] = "SigFox"

		// Push the package to Redis
		err := cache.AppCache.RPush("logs:sigfox-setting-logs", i)
		if err != nil {
			helpers.LogError(err, "Failed to push setting package data log to Redis")
		}

		messageData, err := json.Marshal(i)
		if err != nil {
			helpers.LogError(err, "Failed to serialize parsedData to JSON")
			continue
		}
		mq.AppRabbitMQProducer.SendMessage("event_logs_exchange", "event_logs_queue", string(messageData))
	}
}

// updateDeviceKeepaliveInCacheAndBroadcast updates the keepalive timestamp for a device in the cache and broadcasts changes.
// If the new keepalive timestamp is more recent than the cached one, the cache and relevant logs are updated.
func (h *SigfoxHandler) updateDeviceKeepaliveInCacheAndBroadcast(parsedData map[string]any, deviceID string) error {

	// Extract the list of keepalive packages from the parsed data.
	keepalivePackages, ok := parsedData["keep_alive_packages"].([]map[string]any)
	if !ok {
		return errors.New("invalid or missing keep_alive_packages data")
	}

	// Return early if there are no keepalive packages.
	if len(keepalivePackages) == 0 {
		return nil
	}

	// Retrieve the timestamp from the first keepalive package.
	timestamp, ok := keepalivePackages[0]["timestamp"].(int)
	if !ok {
		return errors.New("timestamp missing or not an integer in keepalive package")
	}

	// Convert the timestamp to a UTC time string.
	timestampTime := time.Unix(int64(timestamp), 0)
	keepaliveAt := timestampTime.UTC().Format("2006-01-02T15:04:05Z")

	// Retrieve cached device data.
	cachedDevice, err := cache.AppCache.GetDevice(deviceID)
	if err != nil {
		helpers.LogError(err, "Error retrieving device from cache")
		return err
	}

	var happenedAt string
	var settingsAt string

	// Check if there is cached data and the new data is more recent.
	if cachedDevice != nil {
		cachedKeepaliveAtStr, ok := cachedDevice["keepalive_at"].(string)
		if !ok {
			helpers.LogError(nil, "Cached keepalive_at is not a string or missing")
			cachedKeepaliveAtStr = "0001-01-01T00:00:00Z" // Default to the earliest possible timestamp
		}
		happenedAt, ok = cachedDevice["happened_at"].(string)
		if !ok {
			return errors.New("cached happened_at is not a string")
		}
		settingsAt, ok = cachedDevice["settings_at"].(string)
		if !ok {
			return errors.New("cached settings_at is not a string")
		}

		cachedKeepaliveAt, err := time.Parse("2006-01-02T15:04:05Z", cachedKeepaliveAtStr)
		if err != nil {
			return fmt.Errorf("error parsing cached keepalive_at time: %v", err)
		}

		newKeepaliveAt, err := time.Parse("2006-01-02T15:04:05Z", keepaliveAt)
		if err != nil {
			return fmt.Errorf("error parsing new keepalive_at time: %v", err)
		}

		// Update only if the new keepalive timestamp is more recent.
		if !newKeepaliveAt.After(cachedKeepaliveAt) {
			helpers.LogInfo("No update needed. Cached keepalive_at is newer or equal.")
			return nil
		}

	} else {
		happenedAt = "0001-01-01T00:00:00Z"
		settingsAt = "0001-01-01T00:00:00Z"
	}

	// --- Update the device cache (e.g., parking:device:<id>)
	err = cache.AppCache.UpdateKeepaliveAt(deviceID, keepaliveAt, happenedAt, settingsAt)
	if err != nil {
		helpers.LogError(err, "Failed to update device keepalive timestamp in cache")
		return err
	}

	// --- Log updates for PostgreSQL synchronization (e.g., logs:device-keepalive-at)
	logPayload := map[string]any{
		"device_id":    deviceID,
		"keepalive_at": keepaliveAt,
	}

	// Push the log entry to Redis for PostgreSQL update processing.
	err = cache.AppCache.RPush("logs:device-keepalive-at", logPayload)
	if err != nil {
		helpers.LogError(err, "Failed to push keepalive update log to Redis")
	}

	// Broadcast the update to clients using Socket.IO.
	app.SocketIO.BroadcastToNamespace("/", "keepalive-event", logPayload)
	helpers.LogInfo("Broadcasted keepalive event for device %s", deviceID)

	return nil
}

// updateDeviceSettingsInCacheAndBroadcast updates the settings timestamp for a device in the cache and broadcasts changes.
// If the new settings timestamp is more recent than the cached one, the cache and relevant logs are updated.
func (h *SigfoxHandler) updateDeviceSettingsInCacheAndBroadcast(parsedData map[string]any, deviceID string) (bool, error) {
	// Extract the list of settings packages from the parsed data.
	settingsPackages, ok := parsedData["settings_packages"].([]map[string]any)

	if !ok {
		return false, errors.New("invalid or missing settings_packages data")
	}

	// Return early if there are no settings packages.
	if len(settingsPackages) == 0 {
		return false, nil
	}

	// Retrieve the timestamp from the first settings package.
	firstSettingsPakage := settingsPackages[0]
	timestamp, ok := firstSettingsPakage["timestamp"].(int)
	if !ok {
		return false, errors.New("timestamp missing or not an integer in settings package")
	}

	// Convert the timestamp to a UTC time string.
	timestampTime := time.Unix(int64(timestamp), 0)
	settingsAt := timestampTime.UTC().Format("2006-01-02T15:04:05Z")

	// Retrieve cached device data.
	cachedDevice, err := cache.AppCache.GetDevice(deviceID)
	if err != nil {
		helpers.LogError(err, "Error retrieving device from cache")
		return false, err
	}

	var happenedAt string
	var keepaliveAt string

	// Check if there is cached data and the new data is more recent.
	if cachedDevice != nil {

		cachedSettingsAtStr, ok := cachedDevice["settings_at"].(string)
		if !ok {
			helpers.LogError(nil, "Cached settings_at is not a string or missing")
			cachedSettingsAtStr = "0001-01-01T00:00:00Z" // Default to the earliest possible timestamp
		}
		happenedAt, ok = cachedDevice["happened_at"].(string)
		if !ok {
			return false, errors.New("cached happened_at is not a string")
		}
		keepaliveAt, ok = cachedDevice["keepalive_at"].(string)
		if !ok {
			return false, errors.New("cached keepalive_at is not a string")
		}

		cachedSettingsAt, err := time.Parse("2006-01-02T15:04:05Z", cachedSettingsAtStr)
		if err != nil {
			return false, fmt.Errorf("error parsing cached settings_at time: %v", err)
		}

		newSettingsAt, err := time.Parse("2006-01-02T15:04:05Z", settingsAt)
		if err != nil {
			return false, fmt.Errorf("error parsing new settings_at time: %v", err)
		}

		// Update only if the new settings timestamp is more recent.
		if !newSettingsAt.After(cachedSettingsAt) {
			helpers.LogInfo("No update needed. Cached settings_at is newer or equal.")
			return false, nil
		}

	} else {
		happenedAt = "0001-01-01T00:00:00Z"
		keepaliveAt = "0001-01-01T00:00:00Z"
	}

	// --- Update the device cache (e.g., parking:device:<id>)
	err = cache.AppCache.UpdateSettingsAt(deviceID, settingsAt, happenedAt, keepaliveAt)
	if err != nil {
		helpers.LogError(err, "Failed to update device settings timestamp in cache")
		return false, err
	}

	// --- Log updates for PostgreSQL synchronization (e.g., logs:device-settings-at)
	logPayload := map[string]any{
		"device_id":   deviceID,
		"settings_at": settingsAt,
	}

	err = cache.AppCache.RPush("logs:device-settings-at", logPayload)
	if err != nil {
		helpers.LogError(err, "Failed to push device settings_at to Redis")
	}

	// Broadcast the update to clients using Socket.IO.
	app.SocketIO.BroadcastToNamespace("/", "settings-event", logPayload)
	helpers.LogInfo("Broadcasted settings event for device %s", deviceID)

	return true, nil
}

// updateDeviceCacheAndBroadcast updates the device data cache and broadcasts changes if the incoming data is newer than what's in the cache.
func (h *SigfoxHandler) updateDeviceCacheAndBroadcast(parsedData map[string]any, deviceId string) error {
	// Extract the list of parking packages from the parsed data
	latestParkingPackage, ok := parsedData["parking_packages"].([]map[string]any)
	if !ok {
		return errors.New("invalid or missing parking_packages data")
	}

	// Return early if there are no parking packages
	if len(latestParkingPackage) == 0 {
		return nil
	}

	// Retrieve the timestamp from the first parking package
	timestamp, ok := latestParkingPackage[0]["timestamp"].(int)
	if !ok {
		return errors.New("timestamp missing or not an integer")
	}

	// Convert the timestamp to a UTC time string
	timestampTime := time.Unix(int64(timestamp), 0)
	happenedAt := timestampTime.UTC().Format("2006-01-02T15:04:05Z")

	// Retrieve cached device data
	cachedDevice, err := cache.AppCache.GetDevice(deviceId)
	if err != nil {
		helpers.LogError(err, "Error retrieving device from cache")
		return err
	}

	// Check if there is cached data and the new data is more recent
	if cachedDevice != nil {
		cachedHappenedAtStr, ok := cachedDevice["happened_at"].(string)
		if !ok {
			return errors.New("cached happened_at is not a string")
		}

		cachedHappenedAt, err := time.Parse("2006-01-02T15:04:05Z", cachedHappenedAtStr)
		if err != nil {
			return fmt.Errorf("error parsing cached happened_at time: %v", err)
		}

		newHappenedAt, err := time.Parse("2006-01-02T15:04:05Z", happenedAt)
		if err != nil {
			return fmt.Errorf("error parsing new happened_at time: %v", err)
		}

		// Proceed with update if the new data is more recent
		if newHappenedAt.After(cachedHappenedAt) {
			return h.processParkingEvent(parsedData, deviceId, happenedAt, latestParkingPackage)
		}

		helpers.LogInfo("No update needed. Cached happened_at is newer or equal.")
		return nil
	}

	// If no cached data exists, process the event as a new entry
	return h.processParkingEvent(parsedData, deviceId, happenedAt, latestParkingPackage)
}

func (h *SigfoxHandler) processParkingEvent(
	parsedData map[string]any,
	deviceId string,
	happenedAt string,
	latestParkingPackage []map[string]any,
) error {
	// Extract the firmware version as a float64
	firmwareVersionFloat, ok := parsedData["firmware_version"].(float64)
	if !ok {
		return errors.New("firmware_version missing or not a float64")
	}

	// Format the firmware version as a string
	firmwareVersion := fmt.Sprintf("%.2f", firmwareVersionFloat)

	// Extract the beacons data from the parking package
	beacons, ok := latestParkingPackage[0]["beacons"].([]map[string]any)
	if !ok {
		return errors.New("beacons missing or not in the expected format")
	}

	// Determine if the parking spot is occupied
	isOccupied := (latestParkingPackage[0]["is_occupied"].(int)) == 1

	// --- Update the device cache (parking:device:<id>)
	err := cache.AppCache.ProcessParkingEventData(deviceId, firmwareVersion, beacons, happenedAt, isOccupied)
	if err != nil {
		helpers.LogError(err, "Failed to update device cache")
		return err
	}

	// --- Log updates for PostgreSQL synchronization (logs:device-update)
	// Create a payload for logging the update
	payload := map[string]any{
		"firmware_version": firmwareVersion,
		"device_id":        deviceId,
		"happened_at":      happenedAt,
		"is_occupied":      isOccupied,
		"beacons":          beacons,
	}

	// Push the log entry to Redis for PostgreSQL update processing
	err = cache.AppCache.RPush("logs:device-update", payload)
	if err != nil {
		helpers.LogError(err, "Failed to push to Redis logs:device-update")
		return err
	}

	// Broadcast the update to clients using socket.io
	app.SocketIO.BroadcastToNamespace("/", "parking-event", payload)

	return nil
}
