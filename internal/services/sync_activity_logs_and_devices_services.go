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
	items, err := s.cache.LRangeAndDelete("logs:nb-activity-logs")
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
	items, err := s.cache.LRangeAndDelete("logs:device-update-logs")
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

		newHappenedAt, err := time.Parse("2006-01-02T15:04:05.000000000Z", itemMap["happened_at"].(string))
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
