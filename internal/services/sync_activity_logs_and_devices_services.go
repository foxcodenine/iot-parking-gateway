package services

import (
	"sort"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

// SyncActivityLogsAndDevices processes and synchronizes activity logs and device updates from Redis to PostgreSQL.
func (s *Service) SyncActivityLogsAndDevices() {

	// Retrieve activity log data from Redis and delete the key.
	items, err := s.cache.LRangeAndDelete("nb-activity-logs")
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

		// If activity logs are successfully inserted, proceed to update devices.
		if err = s.models.Device.BulkUpdateDevices(activityLogs); err != nil {
			s.errorLog.Printf("Failed to update device records: %v", err)
			return
		}

		// Log successful insertion and update.
		s.infoLog.Printf("Successfully inserted %d activity logs and updated device records into PostgreSQL", len(activityLogs))
	}

}
