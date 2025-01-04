package services

import (
	"sort"

	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

func (s *Service) SyncLoraSettingLogs() {
	// Retrieve setting log data from Redis and delete the key.
	items, err := s.cache.LRangeAndDelete("logs:lora-setting-logs")
	if err != nil {
		// Log error if Redis operations fail.
		s.errorLog.Printf("Error retrieving items from Redis: %v", err)
		return
	}

	// Prepare a slice to hold the converted setting log entries.
	loraSettingLogs := make([]models.LoraSettingLog, 0, len(items))

	// Iterate through each item retrieved from Redis.
	for _, item := range items {
		// Attempt to assert the type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			s.errorLog.Println("Invalid item type: expected map[string]any")
			continue
		}

		// Convert the map to an LoraSettingLog struct.
		settingLog, err := models.NewLoraSettingLog(itemMap)
		if err != nil {
			s.errorLog.Printf("Error converting item to LoraSettingLog: %v", err)
			continue
		}

		// Append the successfully created setting log to the slice.
		loraSettingLogs = append(loraSettingLogs, *settingLog)
	}

	// Optionally sort setting logs by the HappenedAt field, if chronological order is important.
	sort.Slice(loraSettingLogs, func(i, j int) bool {
		return loraSettingLogs[i].HappenedAt.Before(loraSettingLogs[j].HappenedAt)
	})

	if len(loraSettingLogs) > 0 {

		err = s.models.LoraSettingLog.BulkInsert(loraSettingLogs)
		if err != nil {
			s.errorLog.Printf("Failed to insert setting logs to PostgreSQL: %v", err)
			return
		}
		// TODO:  to impliment and do better as keepalive
		// err = s.models.LoraDeviceSettings.BulkUpdate(loraSettingLogs)
		// if err != nil {
		// 	s.errorLog.Printf("Failed to update device setting to PostgreSQL: %v", err)
		// 	return
		// }

		s.infoLog.Printf("Successfully inserted and updated %d setting into PostgreSQL", len(loraSettingLogs))
	}
}
