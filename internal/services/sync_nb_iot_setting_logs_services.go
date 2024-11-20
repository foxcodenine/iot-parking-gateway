package services

import (
	"sort"

	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

func (s *Service) SyncNBIoTSettingLogs() {
	// Retrieve setting log data from Redis and delete the key.
	items, err := s.cache.LRangeAndDelete("nb-setting-logs")
	if err != nil {
		// Log error if Redis operations fail.
		s.errorLog.Printf("Error retrieving items from Redis: %v", err)
		return
	}

	// Prepare a slice to hold the converted setting log entries.
	nbIotSettingLogs := make([]models.NbiotSettingLog, 0, len(items))

	// Iterate through each item retrieved from Redis.
	for _, item := range items {
		// Attempt to assert the type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			s.errorLog.Println("Invalid item type: expected map[string]any")
			continue
		}

		// Convert the map to an NbiotSettingLog struct.
		settingLog, err := models.NewNbiotSettingLog(itemMap)
		if err != nil {
			s.errorLog.Printf("Error converting item to NbiotSettingLog: %v", err)
			continue
		}

		// Append the successfully created setting log to the slice.
		nbIotSettingLogs = append(nbIotSettingLogs, *settingLog)
	}

	// Optionally sort setting logs by the HappenedAt field, if chronological order is important.
	sort.Slice(nbIotSettingLogs, func(i, j int) bool {
		return nbIotSettingLogs[i].HappenedAt.Before(nbIotSettingLogs[j].HappenedAt)
	})

	if len(nbIotSettingLogs) > 0 {

		err = s.models.NbiotSettingLog.BulkInsert(nbIotSettingLogs)
		if err != nil {
			s.errorLog.Printf("Failed to insert setting logs to PostgreSQL: %v", err)
			return
		}
		err = s.models.NbiotDeviceSettings.BulkUpdate(nbIotSettingLogs)
		if err != nil {
			s.errorLog.Printf("Failed to update device setting to PostgreSQL: %v", err)
			return
		}

		s.infoLog.Printf("Successfully inserted and updated %d setting into PostgreSQL", len(nbIotSettingLogs))
	}
}
