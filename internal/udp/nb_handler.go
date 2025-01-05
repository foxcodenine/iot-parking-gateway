package udp

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/firmware"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"

	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/google/uuid"
)

// nbMessageHandler processes incoming UDP messages and logs data to Redis.
func (s *UDPServer) nbMessageHandler(conn *net.UDPConn, data []byte, addr *net.UDPAddr) {

	// Prepare initial reply with message type and timestamp
	reply := []string{"0106"}
	hexTimestamp := helpers.GetCurrentTimestampHex()
	reply = append(reply, hexTimestamp)

	// -----------------------------------------------------------------
	// Convert data to a string and trim any newlines
	// hexStr := string(bytes.TrimSpace(data))
	// -----------------------------------------------------------------
	base64Str := string(bytes.TrimSpace(data))
	bufferBase64, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		handleErrorSendResponse(err, "Error decoding base64", conn, addr, reply)
		return
	}
	// Convert the decoded bytes to a hex string
	hexStr := hex.EncodeToString(bufferBase64)
	// -----------------------------------------------------------------

	// Validate minimum hex string length
	if len(hexStr) < 14 {
		handleErrorSendResponse(errors.New("incoming data too short for parsing"), "Invalid message length", conn, addr, reply)
		return
	}

	// Parse firmware version
	firmwareVersionTmp, nextOffset, err := helpers.ParseHexSubstring(hexStr, 0, 1)
	if err != nil {
		handleErrorSendResponse(err, "Failed to parse firmware version", conn, addr, reply)
		return
	}
	// Divide by 10 to convert to float64 and shift decimal place
	firmwareVersion := float64(firmwareVersionTmp) / 10.0

	// Parse device ID
	deviceID, _, err := helpers.ParseHexSubstring(hexStr, nextOffset, 7)
	if err != nil {
		handleErrorSendResponse(err, "Failed to parse device ID", conn, addr, reply)
		return
	}

	// Check if the device ID is already in the Bloom Filter
	deviceIdentifierKey := fmt.Sprintf("NB-IoT %d", deviceID)
	isDeviceRegistered, err := s.cache.CheckItemInBloomFilter("registered-devices", deviceIdentifierKey)
	if err != nil {
		helpers.LogError(err, "Failed to check Bloom Filter for device ID")
	}

	// If the device ID is not registered, track it for registration and prevent duplicates.
	if !isDeviceRegistered {
		// Add the device to a Redis set for tracking devices that need registration.
		deviceDataKey := fmt.Sprintf("%s %f", deviceIdentifierKey, firmwareVersion)
		if err := s.cache.SAdd("to-register-devices", deviceDataKey); err != nil {
			helpers.LogError(err, "Failed to add device ID to the 'to-register-devices' set")
		}

		// Add the device ID to the Bloom Filter to prevent duplicate registrations in the future.
		if _, err := s.cache.AddItemToBloomFilter("registered-devices", deviceIdentifierKey); err != nil {
			helpers.LogError(err, "Failed to add device ID to the 'registered-devices' Bloom Filter")
		}

	}

	// If the device ID is registered, check if it is soft delete, white listed or black listed.
	if isDeviceRegistered {

		// Retrieve the device data from the cache
		deviceData, err := s.cache.GetDevice(fmt.Sprintf("%d", deviceID))
		if err != nil {
			helpers.LogError(err, "Failed to retrieve device data from cache.")
			sendResponse(conn, addr, reply)
			return
		}

		// Check if the device is soft deleted
		if deletedAt, exists := deviceData["deleted_at"]; exists && deletedAt != nil && deletedAt != "0001-01-01T00:00:00Z" {
			helpers.LogInfo("Device %d is marked as soft deleted. Request ignored.", deviceID)
			sendResponse(conn, addr, reply)
			return
		}

		// Retrieve application settings for device access mode
		deviceAccessMode, err := s.cache.HGet("app:settings", "device_access_mode")
		if err != nil {
			helpers.LogError(err, "Failed to retrieve 'device_access_mode' from application settings.")
			sendResponse(conn, addr, reply)
			return
		}

		// Check access based on whitelist mode
		if deviceAccessMode == "white_list" {
			if isAllowed, ok := deviceData["is_allowed"].(bool); ok && !isAllowed {
				helpers.LogInfo("Device %d is not marked allowed. Request ignored.", deviceID)
				sendResponse(conn, addr, reply)
				return
			}
		}

		// Check access based on blacklist mode
		if deviceAccessMode == "black_list" {
			if isBlocked, ok := deviceData["is_blocked"].(bool); ok && isBlocked {
				helpers.LogInfo("Device %d is marked blocked. Request ignored.", deviceID)
				sendResponse(conn, addr, reply)
				return
			}
		}
	}

	// Generate a new UUID for the RawDataLog entry
	rawUUID, err := uuid.NewV7()
	if err != nil {
		handleErrorSendResponse(err, "Failed to generate UUID for RawDataLog entry", conn, addr, reply)
		return
	}

	// Create a new RawDataLog object to store in Redis.
	rawDataLog := models.RawDataLog{
		ID:              rawUUID,
		DeviceID:        strconv.Itoa(deviceID),
		FirmwareVersion: firmwareVersion,
		NetworkType:     "NB-IoT",
		RawData:         hexStr,
		CreatedAt:       time.Now(),
	}

	// Push the raw data log entry to Redis
	err = s.cache.RPush("logs:raw-data-logs", rawDataLog)
	if err != nil {
		handleErrorSendResponse(err, "Failed to push raw data log to Redis", conn, addr, reply)
		return
	}

	// Debug output for parsed values
	helpers.LogInfo("Network: NB-IoT, Firmware: %.2f, Device ID: %d", firmwareVersion, deviceID)

	// Process firmware-specific data parsing based on the firmware version.
	var parsedData map[string]any
	switch firmwareVersion {
	case 5.3:
		parsedData, err = firmware.NB_53(hexStr)
	case 5.8, 5.9:
		parsedData, err = firmware.NB_58(hexStr)

	default:
		// Send a default response if the firmware version is not handled.
		sendResponse(conn, addr, reply)
		return
	}

	if err != nil {
		handleErrorSendResponse(err, fmt.Sprintf("Failed to parse data from NB_%.0f firmware", firmwareVersion*10), conn, addr, reply)
		return
	}

	// -----------------------------------------------------------------

	// Attempt to update device cache and broadcast the changes.
	err = s.updateDeviceCacheAndBroadcast(parsedData)

	// Check for errors in the update process.
	if err != nil {
		// Log the error with additional context for better troubleshooting.
		helpers.LogError(err, "Failed to update device cache and broadcast changes")
	}

	// Attempt to update device keepalive_at in cache and broadcast the changes.
	err = s.updateDeviceKeepaliveInCacheAndBroadcast(parsedData)

	// Check for errors in the update process.
	if err != nil {
		// Log the error with additional context for better troubleshooting.
		helpers.LogError(err, "Failed to update device keepalive_at in cache and broadcast it")
	}

	// Attempt to update device settings_at in cache, check if device_settings should be updated and broadcast settings_at.
	updateDeviceSettings, err := s.updateDeviceSettingsInCacheAndBroadcast(parsedData)

	// Check for errors in the update process.
	if err != nil {
		// Log the error with additional context for better troubleshooting.
		helpers.LogError(err, "Failed to update device settings_at in cache and broadcast it")
	}

	// Push parsed parking data packages to Redis.
	for _, i := range parsedData["parking_packages"].([]map[string]any) {

		i["firmware_version"] = parsedData["firmware_version"]
		i["device_id"] = fmt.Sprintf("%d", parsedData["device_id"])
		i["raw_id"] = rawUUID
		i["event_id"] = 26
		i["network_type"] = "NB-IoT"

		err := s.cache.RPush("logs:activity-logs", i)
		if err != nil {
			helpers.LogError(err, "Failed to push parking package data log to Redis")
		}

		messageData, err := json.Marshal(i)
		if err != nil {
			helpers.LogError(err, "Failed to serialize parsedData to JSON")
			continue
		}
		s.mqProducer.SendMessage("event_logs_exchange", "event_logs_queue", string(messageData))
	}

	// Push parsed keepalive data to Redis.
	for _, i := range parsedData["keep_alive_packages"].([]map[string]any) {
		i["firmware_version"] = parsedData["firmware_version"]
		i["device_id"] = fmt.Sprintf("%d", parsedData["device_id"])
		i["raw_id"] = rawUUID
		i["event_id"] = 6
		i["network_type"] = "NB-IoT"

		err := s.cache.RPush("logs:nb-keepalive-logs", i)
		if err != nil {
			helpers.LogError(err, "Failed to push keepalive package data log to Redis")
		}

		messageData, err := json.Marshal(i)
		if err != nil {
			helpers.LogError(err, "Failed to serialize parsedData to JSON")
			continue
		}
		s.mqProducer.SendMessage("event_logs_exchange", "event_logs_queue", string(messageData))
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
		i["device_id"] = fmt.Sprintf("%d", parsedData["device_id"])
		i["raw_id"] = rawUUID
		i["event_id"] = 25 // Assuming 25 is the event ID for setting logs
		i["network_type"] = "NB-IoT"

		// Push the package to Redis
		err := s.cache.RPush("logs:nb-setting-logs", i)
		if err != nil {
			helpers.LogError(err, "Failed to push setting package data log to Redis")
		}

		messageData, err := json.Marshal(i)
		if err != nil {
			helpers.LogError(err, "Failed to serialize parsedData to JSON")
			continue
		}
		s.mqProducer.SendMessage("event_logs_exchange", "event_logs_queue", string(messageData))
	}

	// time.Sleep(1 * time.Second)
	// s.services.RegisterNewDevices()
	// s.services.SyncActivityLogs()
	// s.services.SyncNBIoTKeepaliveLogs()
	// s.services.SyncNBIoTSettingLogs()

	// Send the final response back to the UDP client to confirm processing.
	sendResponse(conn, addr, reply)
}

// sendResponse sends a structured reply back to the UDP client.
func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, reply []string) {
	response := []byte(strings.Join(reply, "") + "\n")
	_, err := conn.WriteToUDP(response, addr)
	if err != nil {
		helpers.LogError(err, "Failed to send response to client")
	}
}

// handleErrorSendResponse logs an error, sends a response to the client, and exits the handler.
func handleErrorSendResponse(err error, message string, conn *net.UDPConn, addr *net.UDPAddr, reply []string) {
	helpers.LogError(err, message, 3)
	sendResponse(conn, addr, reply)
}

// updateDeviceKeepaliveInCacheAndBroadcast updates the keepalive timestamp for a device in the cache and broadcasts changes.
// If the new keepalive timestamp is more recent than the cached one, the cache and relevant logs are updated.
func (s *UDPServer) updateDeviceKeepaliveInCacheAndBroadcast(parsedData map[string]any) error {
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
	deviceID := fmt.Sprintf("%d", parsedData["device_id"])

	// Retrieve cached device data.
	cachedDevice, err := s.cache.GetDevice(deviceID)
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
	err = s.cache.UpdateKeepaliveAt(deviceID, keepaliveAt, happenedAt, settingsAt)
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
	err = s.cache.RPush("logs:device-keepalive-at", logPayload)
	if err != nil {
		helpers.LogError(err, "Failed to push device keepalive_at to Redis")
	}

	// Broadcast the update to clients using Socket.IO.
	s.SocketIO.BroadcastToNamespace("/", "keepalive-event", logPayload)
	helpers.LogInfo("Broadcasted keepalive event for device %s", deviceID)
	// TODO: implement "keepalive-event" in frontend

	return nil
}

// updateDeviceSettingsInCacheAndBroadcast updates the settings timestamp for a device in the cache and broadcasts changes.
// If the new settings timestamp is more recent than the cached one, the cache and relevant logs are updated.
func (s *UDPServer) updateDeviceSettingsInCacheAndBroadcast(parsedData map[string]any) (bool, error) {
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
	deviceID := fmt.Sprintf("%d", parsedData["device_id"])

	// Retrieve cached device data.
	cachedDevice, err := s.cache.GetDevice(deviceID)
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
	err = s.cache.UpdateSettingsAt(deviceID, settingsAt, happenedAt, keepaliveAt)
	if err != nil {
		helpers.LogError(err, "Failed to update device settings timestamp in cache")
		return false, err
	}

	// --- Log updates for PostgreSQL synchronization (e.g., logs:device-settings-at)
	logPayload := map[string]any{
		"device_id":   deviceID,
		"settings_at": settingsAt,
	}

	// Push the log entry to Redis for PostgreSQL update processing.
	err = s.cache.RPush("logs:device-settings-at", logPayload)
	if err != nil {
		helpers.LogError(err, "Failed to push device settings_at to Redis")
	}

	// Broadcast the update to clients using Socket.IO.
	s.SocketIO.BroadcastToNamespace("/", "settings-event", logPayload)
	helpers.LogInfo("Broadcasted settings event for device %s", deviceID)
	// TODO: implement "settings-event" in frontend

	return true, nil
}

// updateDeviceCacheAndBroadcast updates the device data cache and broadcasts changes if the incoming data is newer than what's in the cache.
// updateDeviceCacheAndBroadcast updates the device data cache and broadcasts changes if the incoming data is newer than what's in the cache.
func (s *UDPServer) updateDeviceCacheAndBroadcast(parsedData map[string]any) error {
	// Extract the list of parking packages from the parsed data.
	latestParkingPackage, ok := parsedData["parking_packages"].([]map[string]any)
	if !ok {
		return errors.New("invalid or missing parking_packages data")
	}

	// Return early if there are no parking packages.
	if len(latestParkingPackage) == 0 {
		return nil
	}

	// Retrieve the timestamp from the first parking package.
	timestamp, ok := latestParkingPackage[0]["timestamp"].(int)
	if !ok {
		return errors.New("timestamp missing or not an integer")
	}

	// Convert the timestamp to a UTC time string.
	timestampTime := time.Unix(int64(timestamp), 0)
	happenedAt := timestampTime.UTC().Format("2006-01-02T15:04:05Z")
	deviceId := fmt.Sprintf("%d", parsedData["device_id"])

	// Retrieve cached device data.
	cachedDevice, err := s.cache.GetDevice(deviceId)
	if err != nil {
		helpers.LogError(err, "Error retrieving device from cache")
		return err
	}

	// Check if there is cached data and the new data is more recent.
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

		// Proceed with update if the new data is more recent.
		if newHappenedAt.After(cachedHappenedAt) {
			return s.processParkingEvent(parsedData, happenedAt, latestParkingPackage)
		}

		helpers.LogInfo("No update needed. Cached happened_at is newer or equal.")
		return nil
	}

	// If no cached data exists, process the event as a new entry.
	return s.processParkingEvent(parsedData, happenedAt, latestParkingPackage)
}

// processParkingEvent processes a parking event by updating the cache, logging to Redis, and broadcasting changes.
func (s *UDPServer) processParkingEvent(parsedData map[string]any, happenedAt string, latestParkingPackage []map[string]any) error {
	// Extract the firmware version as a float64.
	firmwareVersionFloat, ok := parsedData["firmware_version"].(float64)
	if !ok {
		return errors.New("firmware_version missing or not a float64")
	}

	// Format the firmware version as a string.
	firmwareVersion := fmt.Sprintf("%.2f", firmwareVersionFloat)

	// Extract the beacons data from the parking package.
	beacons, ok := latestParkingPackage[0]["beacons"].([]map[string]any)
	if !ok {
		return errors.New("beacons missing or not in the expected format")
	}

	// Determine if the parking spot is occupied.
	isOccupied := (latestParkingPackage[0]["is_occupied"].(int)) == 1

	// --- Update the device cache (parking:device:<id>)
	err := s.cache.ProcessParkingEventData(fmt.Sprintf("%d", parsedData["device_id"]), firmwareVersion, beacons, happenedAt, isOccupied)
	if err != nil {
		helpers.LogError(err, "Failed to update device cache")
		return err
	}

	// --- Log updates for PostgreSQL synchronization (logs:device-update).
	payload := map[string]any{
		"firmware_version": firmwareVersion,
		"device_id":        fmt.Sprintf("%d", parsedData["device_id"]),
		"happened_at":      happenedAt,
		"is_occupied":      isOccupied,
		"beacons":          beacons,
	}

	// Push the log entry to Redis for PostgreSQL update processing.
	err = s.cache.RPush("logs:device-update", payload)
	if err != nil {
		helpers.LogError(err, "Failed to push to Redis logs:device-update")
	}

	// Broadcast the update to clients using Socket.IO.
	s.SocketIO.BroadcastToNamespace("/", "parking-event", payload)

	return nil
}
