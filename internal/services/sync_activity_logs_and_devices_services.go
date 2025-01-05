package services

import (
	"database/sql"
	"sort"
	"strconv"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

// SyncActivityLogs processes and synchronizes activity logs and device updates from Redis to PostgreSQL.
func (s *Service) SyncActivityLogs() {

	// Retrieve activity log data from Redis and delete the key.
	items, err := s.cache.LRangeAndDelete("logs:activity-logs")
	if err != nil {
		// Log error if Redis operations fail.
		s.errorLog.Printf("Error retrieving items from Redis: %v", err)
		return
	}

	// Prepare a slice to hold the converted activity log entries.
	activityLogs := make([]models.ActivityLog, 0, len(items))

	// Iterate through each item retrieved from Redis.
	for _, item := range items {

		// Attempt to assert the type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			s.errorLog.Println("Invalid item type: expected map[string]any")
			continue
		}

		// Convert the map to an ActivityLog struct.
		activityLog, err := models.NewActivityLog(itemMap)
		if err != nil {
			helpers.LogError(err, "")
		}

		// Append the successfully created activity log to the slice.
		activityLogs = append(activityLogs, *activityLog)

	}

	// Sort activity logs by the HappenedAt field.
	sort.Slice(activityLogs, func(i, j int) bool {
		return activityLogs[i].HappenedAt.Before(activityLogs[j].HappenedAt)
	})

	if len(activityLogs) > 0 {

		// Attempt to bulk insert all activity logs into PostgreSQL.
		if err = s.models.ActivityLog.BulkInsert(activityLogs); err != nil {
			s.errorLog.Printf("Failed to insert activity logs to PostgreSQL: %v", err)
			return
		}

		// Log successful insertion and update.
		s.infoLog.Printf("Successfully inserted %d activity logs records into PostgreSQL", len(activityLogs))
	}

}

// SyncDevices processes and synchronizes device updates from Redis to PostgreSQL.
func (s *Service) SyncDevices() {

	// Retrieve activity log data from Redis and delete the key.
	items, err := s.cache.LRangeAndDelete("logs:device-update")
	if err != nil {
		// Log error if Redis operations fail.
		s.errorLog.Printf("Error retrieving items from Redis: %v", err)
		return
	}

	// Prepare a slice to hold the converted activity log entries.
	deviceUpdateLogs := make([]models.Device, 0, len(items))

	// Iterate through each item retrieved from Redis.
	for _, item := range items {

		// Attempt to assert the type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			s.errorLog.Println("Invalid item type: expected map[string]any")
			continue
		}

		firmwareVersionStr, ok := itemMap["firmware_version"].(string)
		if !ok {
			helpers.LogError(nil, "firmware_version is not a string")
			continue
		}
		firmwareVersion, err := strconv.ParseFloat(firmwareVersionStr, 64)
		if err != nil {
			helpers.LogError(nil, "error converting firmware_version to float64:")
			continue
		}

		newHappenedAt, err := time.Parse("2006-01-02T15:04:05Z", itemMap["happened_at"].(string))
		if err != nil {
			helpers.LogError(err, "error parsing new happened_at time")
			continue
		}

		device := models.Device{
			DeviceID:        itemMap["device_id"].(string),
			FirmwareVersion: firmwareVersion,
			IsOccupied:      itemMap["is_occupied"].(bool),
			HappenedAt:      newHappenedAt,
		}
		if itemMap["beacons_json"] != nil {
			device.BeaconsJSON = itemMap["beacons_json"].(sql.NullString)
			device.ParseBeaconsJSON()
		}

		// Append the successfully created activity log to the slice.
		deviceUpdateLogs = append(deviceUpdateLogs, device)

	}

	// Sort activity logs by the HappenedAt field.
	sort.Slice(deviceUpdateLogs, func(i, j int) bool {
		return deviceUpdateLogs[i].HappenedAt.Before(deviceUpdateLogs[j].HappenedAt)
	})

	if len(deviceUpdateLogs) > 0 {

		// If activity logs are successfully inserted, proceed to update devices.
		if err = s.models.Device.BulkUpdateDevices(deviceUpdateLogs); err != nil {
			s.errorLog.Printf("Failed to update device records: %v", err)
			return
		}

		s.infoLog.Printf("Successfully updated %d device records into PostgreSQL", len(deviceUpdateLogs))
	}

}

// SyncDevicesKeepaliveAt processes and synchronizes device keepalive updates from Redis to PostgreSQL.
func (s *Service) SyncDevicesKeepaliveAt() {

	// Retrieve keepalive log data from Redis and delete the key.
	items, err := s.cache.LRangeAndDelete("logs:device-keepalive-at")
	if err != nil {
		// Log an error if Redis operations fail and return early.
		s.errorLog.Printf("Error retrieving items from Redis: %v", err)
		return
	}

	// Prepare a slice to hold the converted keepalive log entries.
	deviceUpdateLogs := make([]models.Device, 0, len(items))

	// Iterate through each item retrieved from Redis.
	for _, item := range items {

		// Attempt to assert the item type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			s.errorLog.Println("Invalid item type: expected map[string]any")
			continue
		}

		// Ensure the `device_id` field exists and is a string.
		deviceID, ok := itemMap["device_id"].(string)
		if !ok || deviceID == "" {
			s.errorLog.Println("Missing or invalid device_id in item map")
			continue
		}

		// Ensure the `keepalive_at` field exists and is a valid timestamp.
		keepaliveAtStr, ok := itemMap["keepalive_at"].(string)
		if !ok || keepaliveAtStr == "" {
			s.errorLog.Printf("Missing or invalid keepalive_at for device %s", deviceID)
			continue
		}

		// Parse the `keepalive_at` timestamp.
		newKeepaliveAt, err := time.Parse("2006-01-02T15:04:05Z", keepaliveAtStr)
		if err != nil {
			helpers.LogError(err, "Error parsing keepalive_at timestamp")
			continue
		}

		// Create a Device struct for this keepalive log entry.
		device := models.Device{
			DeviceID:    deviceID,
			KeepaliveAt: newKeepaliveAt,
		}

		// Append the device to the slice for batch processing.
		deviceUpdateLogs = append(deviceUpdateLogs, device)
	}

	// Sort the keepalive logs by the KeepaliveAt field.
	sort.Slice(deviceUpdateLogs, func(i, j int) bool {
		return deviceUpdateLogs[i].KeepaliveAt.Before(deviceUpdateLogs[j].KeepaliveAt)
	})

	// If there are no valid entries, return early.
	if len(deviceUpdateLogs) == 0 {
		// s.infoLog.Println("No valid keepalive logs to process")
		return
	}

	// Attempt to bulk update the keepalive timestamps in PostgreSQL.
	if err := s.models.Device.BulkUpdateDevicesKeepalive(deviceUpdateLogs); err != nil {
		s.errorLog.Printf("Failed to bulk update keepalive timestamps for devices: %v", err)
		return
	}

	// Log the successful update.
	s.infoLog.Printf("Successfully updated keepalive timestamps for %d devices in PostgreSQL", len(deviceUpdateLogs))
}

// SyncDevicesSettings processes and synchronizes device settings updates from Redis to PostgreSQL.
func (s *Service) SyncDevicesSettingsAt() {

	// Retrieve settings log data from Redis and delete the key.
	items, err := s.cache.LRangeAndDelete("logs:device-settings-at")
	if err != nil {
		// Log an error if Redis operations fail and return early.
		s.errorLog.Printf("Error retrieving items from Redis: %v", err)
		return
	}

	// Prepare a slice to hold the converted devices entries.
	deviceUpdateLogs := make([]models.Device, 0, len(items))

	// Iterate through each item retrieved from Redis.
	for _, item := range items {

		// Attempt to assert the item type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			s.errorLog.Println("Invalid item type: expected map[string]any")
			continue
		}

		// Ensure the `device_id` field exists and is a string.
		deviceID, ok := itemMap["device_id"].(string)
		if !ok || deviceID == "" {
			s.errorLog.Println("Missing or invalid device_id in item map")
			continue
		}

		// Ensure the `settings_at` field exists and is a valid timestamp.
		settingsAtStr, ok := itemMap["settings_at"].(string)
		if !ok || settingsAtStr == "" {
			s.errorLog.Printf("Missing or invalid settings_at for device %s", deviceID)
			continue
		}

		// Parse the `settings_at` timestamp.
		newSettingsAt, err := time.Parse("2006-01-02T15:04:05Z", settingsAtStr)
		if err != nil {
			helpers.LogError(err, "Error parsing settings_at timestamp")
			continue
		}

		// Create a Device struct for this settings log entry.
		device := models.Device{
			DeviceID:   deviceID,
			SettingsAt: newSettingsAt,
		}

		// Append the device to the slice for batch processing.
		deviceUpdateLogs = append(deviceUpdateLogs, device)
	}

	// Sort the settings logs by the SettingsAt field.
	sort.Slice(deviceUpdateLogs, func(i, j int) bool {
		return deviceUpdateLogs[i].SettingsAt.Before(deviceUpdateLogs[j].SettingsAt)
	})

	// If there are no valid entries, return early.
	if len(deviceUpdateLogs) == 0 {
		// s.infoLog.Println("No valid Settings logs to process")
		return
	}

	// Attempt to bulk update the Settings timestamps in PostgreSQL.
	if err := s.models.Device.BulkUpdateDevicesSettings(deviceUpdateLogs); err != nil {
		s.errorLog.Printf("Failed to bulk update settings timestamps for devices: %v", err)
		return
	}

	// Log the successful update.
	s.infoLog.Printf("Successfully updated settings timestamps for %d devices in PostgreSQL", len(deviceUpdateLogs))
}
