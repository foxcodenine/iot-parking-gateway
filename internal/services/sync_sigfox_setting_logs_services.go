package services

import (
	"sort"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

func (s *Service) SyncSigfoxSettingLogs() {
	// Retrieve setting log data from Redis and delete the key.
	items, err := s.cache.LRangeAndDelete("logs:sigfox-setting-logs")
	if err != nil {
		// Log error if Redis operations fail.
		helpers.LogError(err, "Error retrieving sigfox-setting-logs from Redis")
		return
	}

	// Process items and convert them to sigfoxSettingLogs.
	sigfoxSettingLogs, sigfoxDeviceSetting, err := s.ProcessSigfoxSettingLogs(items)
	if err != nil {
		helpers.LogError(err, "Error processing NB-IoT setting logs")
	}

	// Track errors independently for the two operations.
	var settingLogsError, deviceSettingsError error

	// Insert logs into the database if there are any.
	if len(sigfoxSettingLogs) > 0 {
		settingLogsError = s.models.SigfoxSettingLog.BulkInsert(sigfoxSettingLogs)
		if settingLogsError != nil {
			helpers.LogError(settingLogsError, "Failed to insert sigfox_settings_logs to PostgreSQL")
		} else {
			helpers.LogInfo("Successfully inserted %d sigfox_settings_logs in PostgreSQL", len(sigfoxSettingLogs))
		}
	}

	// Update device settings based on the logs if there are any.
	if len(sigfoxDeviceSetting) > 0 {
		deviceSettingsError = s.models.SigfoxDeviceSettings.BulkUpdate(sigfoxDeviceSetting)
		if deviceSettingsError != nil {
			helpers.LogError(deviceSettingsError, "Failed to update sigfox_device_settings in PostgreSQL")
		} else {
			helpers.LogInfo("Successfully updated %d sigfox_device_settings in PostgreSQL", len(sigfoxDeviceSetting))
		}
	}

	// If both operations failed, log an overarching error.
	if settingLogsError != nil && deviceSettingsError != nil {
		helpers.LogError(nil, "Both sigfox_settings_logs and sigfox_device_settings operations failed")
	}
}

// ProcessSigfoxSettingLogs processes raw Redis items into a slice of SigfoxSettingLog models.
func (s *Service) ProcessSigfoxSettingLogs(items []any) ([]models.SigfoxSettingLog, []models.SigfoxSettingLog, error) {
	// Prepare a slice to hold the converted setting log entries.
	sigfoxSettingLogs := make([]models.SigfoxSettingLog, 0, len(items))

	var sigfoxDeviceSetting []models.SigfoxSettingLog

	for _, item := range items {
		// Assert the item type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			// Log the error but continue processing other items.
			helpers.LogInfo("Invalid item type: expected map[string]any")
			continue
		}

		// Convert the map to an SigfoxSettingLog struct.
		settingLog, err := models.NewSigfoxSettingLog(itemMap)
		if err != nil {
			// Log the conversion error but continue processing other items.
			helpers.LogError(err, "Error converting item to SigfoxSettingLog")
			continue
		}

		// Append the successfully created setting log to the slice.
		sigfoxSettingLogs = append(sigfoxSettingLogs, *settingLog)

		if itemMap["update_device_settings"] == true {
			sigfoxDeviceSetting = append(sigfoxDeviceSetting, *settingLog)
		}
	}

	// Optionally sort the logs by the HappenedAt field for chronological order.
	sort.Slice(sigfoxSettingLogs, func(i, j int) bool {
		return sigfoxSettingLogs[i].HappenedAt.Before(sigfoxSettingLogs[j].HappenedAt)
	})

	return sigfoxSettingLogs, sigfoxDeviceSetting, nil
}
