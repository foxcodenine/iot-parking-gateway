package services

import (
	"sort"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

func (s *Service) SyncNBIoTSettingLogs() {
	// Retrieve setting log data from Redis and delete the key.
	items, err := s.cache.LRangeAndDelete("logs:nb-setting-logs")
	if err != nil {
		// Log error if Redis operations fail.
		helpers.LogError(err, "Error retrieving nb-setting-logs from Redis")
		return
	}

	// Process items and convert them to NbiotSettingLogs.
	nbIotSettingLogs, nbIotDeviceSetting, err := s.ProcessNBIoTSettingLogs(items)
	if err != nil {
		helpers.LogError(err, "Error processing NB-IoT setting logs")
	}

	// Track errors independently for the two operations.
	var settingLogsError, deviceSettingsError error

	// Insert logs into the database if there are any.
	if len(nbIotSettingLogs) > 0 {
		settingLogsError = s.models.NbiotSettingLog.BulkInsert(nbIotSettingLogs)
		if settingLogsError != nil {
			helpers.LogError(settingLogsError, "Failed to insert nbiot_settings_logs to PostgreSQL")
		} else {
			helpers.LogInfo("Successfully inserted %d nbiot_settings_logs in PostgreSQL", len(nbIotSettingLogs))
		}
	}

	// Update device settings based on the logs if there are any.
	if len(nbIotDeviceSetting) > 0 {
		deviceSettingsError = s.models.NbiotDeviceSettings.BulkUpdate(nbIotDeviceSetting)
		if deviceSettingsError != nil {
			helpers.LogError(deviceSettingsError, "Failed to update nbiot_device_settings in PostgreSQL")
		} else {
			helpers.LogInfo("Successfully updated %d nbiot_device_settings in PostgreSQL", len(nbIotDeviceSetting))
		}
	}

	// If both operations failed, log an overarching error.
	if settingLogsError != nil && deviceSettingsError != nil {
		helpers.LogError(nil, "Both nbiot_settings_logs and nbiot_device_settings operations failed")
	}
}

// ProcessNBIoTSettingLogs processes raw Redis items into a slice of NbiotSettingLog models.
func (s *Service) ProcessNBIoTSettingLogs(items []any) ([]models.NbiotSettingLog, []models.NbiotSettingLog, error) {
	// Prepare a slice to hold the converted setting log entries.
	nbIotSettingLogs := make([]models.NbiotSettingLog, 0, len(items))

	var nbIotDeviceSetting []models.NbiotSettingLog

	for _, item := range items {
		// Assert the item type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			// Log the error but continue processing other items.
			helpers.LogError(nil, "Invalid item type: expected map[string]any")
			continue
		}

		// Convert the map to an NbiotSettingLog struct.
		settingLog, err := models.NewNbiotSettingLog(itemMap)
		if err != nil {
			// Log the conversion error but continue processing other items.
			helpers.LogError(err, "Error converting item to NbiotSettingLog")
			continue
		}

		// Append the successfully created setting log to the slice.
		nbIotSettingLogs = append(nbIotSettingLogs, *settingLog)

		if itemMap["update_device_settings"] == true {
			nbIotDeviceSetting = append(nbIotDeviceSetting, *settingLog)
		}
	}

	// Optionally sort the logs by the HappenedAt field for chronological order.
	sort.Slice(nbIotSettingLogs, func(i, j int) bool {
		return nbIotSettingLogs[i].HappenedAt.Before(nbIotSettingLogs[j].HappenedAt)
	})

	return nbIotSettingLogs, nbIotDeviceSetting, nil
}
