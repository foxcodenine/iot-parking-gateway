package services

import (
	"sort"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

func (s *Service) SyncLoraSettingLogs() {
	// Retrieve setting log data from Redis and delete the key.
	items, err := s.cache.LRangeAndDelete("logs:lora-setting-logs")
	if err != nil {
		// Log error if Redis operations fail.
		helpers.LogError(err, "Error retrieving lora-setting-logs from Redis")
		return
	}

	// Process items and convert them to loraSettingLogs.
	loraSettingLogs, loraDeviceSetting, err := s.ProcessLoraSettingLogs(items)
	if err != nil {
		helpers.LogError(err, "Error processing NB-IoT setting logs")
	}

	// Track errors independently for the two operations.
	var settingLogsError, deviceSettingsError error

	// Insert logs into the database if there are any.
	if len(loraSettingLogs) > 0 {
		settingLogsError = s.models.LoraSettingLog.BulkInsert(loraSettingLogs)
		if settingLogsError != nil {
			helpers.LogError(settingLogsError, "Failed to insert lora_settings_logs to PostgreSQL")
		} else {
			helpers.LogInfo("Successfully inserted %d lora_settings_logs in PostgreSQL", len(loraSettingLogs))
		}
	}

	// Update device settings based on the logs if there are any.
	if len(loraDeviceSetting) > 0 {
		deviceSettingsError = s.models.LoraDeviceSettings.BulkUpdate(loraDeviceSetting)
		if deviceSettingsError != nil {
			helpers.LogError(err, "Failed to update lora_device_settings in PostgreSQL")
		} else {
			helpers.LogInfo("Successfully updated %d lora_device_settings in PostgreSQL", len(loraDeviceSetting))
		}
	}

	// If both operations failed, log an overarching error.
	if settingLogsError != nil && deviceSettingsError != nil {
		helpers.LogError(nil, "Both lora_settings_logs and lora_device_settings operations failed")
	}
}

// ProcessLoraSettingLogs processes raw Redis items into a slice of LoraSettingLog models.
func (s *Service) ProcessLoraSettingLogs(items []any) ([]models.LoraSettingLog, []models.LoraSettingLog, error) {
	// Prepare a slice to hold the converted setting log entries.
	loraSettingLogs := make([]models.LoraSettingLog, 0, len(items))

	var loraDeviceSetting []models.LoraSettingLog

	for _, item := range items {
		// Assert the item type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			// Log the error but continue processing other items.
			helpers.LogError(nil, "Invalid item type: expected map[string]any")
			continue
		}

		// Convert the map to an LoraSettingLog struct.
		settingLog, err := models.NewLoraSettingLog(itemMap)
		if err != nil {
			// Log the conversion error but continue processing other items.
			helpers.LogError(err, "Error converting item to LoraSettingLog")
			continue
		}

		// Append the successfully created setting log to the slice.
		loraSettingLogs = append(loraSettingLogs, *settingLog)

		if itemMap["update_device_settings"] == true {
			loraDeviceSetting = append(loraDeviceSetting, *settingLog)
		}
	}

	// Optionally sort the logs by the HappenedAt field for chronological order.
	sort.Slice(loraSettingLogs, func(i, j int) bool {
		return loraSettingLogs[i].HappenedAt.Before(loraSettingLogs[j].HappenedAt)
	})

	return loraSettingLogs, loraDeviceSetting, nil
}
