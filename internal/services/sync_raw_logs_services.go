package services

import (
	"fmt"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/google/uuid"
)

// SyncRawLogs processes and synchronizes raw logs from Redis to PostgreSQL.
func (s *Service) SyncRawLogs() {

	items, err := s.cache.LRangeAndDelete("logs:raw-data-logs")
	if err != nil {
		helpers.LogError(err, "Error retrieving items from Redis")
		return
	}

	var rawDataLogs []models.RawDataLog

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			helpers.LogInfo("Invalid item type: expected map[string]interface{}")
			continue
		}

		// Convert id field from string to uuid.UUID
		uuidStr, ok := itemMap["id"].(string)
		if !ok {
			helpers.LogInfo("Invalid UUID format: expected string")
			continue
		}

		uuidValue, err := uuid.Parse(uuidStr)
		if err != nil {
			helpers.LogError(err, fmt.Sprintf("Failed to parse UUID %s", uuidStr))
			continue
		}

		// Parse created_at field from string to time.Time
		createAtStr, ok := itemMap["created_at"].(string)
		if !ok {
			helpers.LogInfo("Invalid created_at format: expected string")
			continue
		}

		createdAt, err := time.Parse(time.RFC3339, createAtStr)
		if err != nil {
			helpers.LogError(err, fmt.Sprintf("Failed to parse created_at %s\n", createAtStr))
			continue
		}

		// Check for and convert other fields safely
		deviceID, ok := itemMap["device_id"].(string)
		if !ok {
			helpers.LogInfo("Invalid device_id format: expected string")
			continue
		}

		firmwareVersion, ok := itemMap["firmware_version"].(float64) // Redis stores numbers as float64
		if !ok {
			helpers.LogInfo("Invalid firmware_version format: expected float64")
			continue
		}

		networkType, ok := itemMap["network_type"].(string)
		if !ok {
			helpers.LogInfo("Invalid network_type format: expected string")
			continue
		}

		rawData, ok := itemMap["raw_data"].(string)
		if !ok {
			helpers.LogInfo("Invalid raw_data format: expected string")
			continue
		}

		rawLog := models.RawDataLog{
			ID:              uuidValue,
			DeviceID:        deviceID,
			FirmwareVersion: firmwareVersion,
			NetworkType:     networkType,
			RawData:         rawData,
			CreatedAt:       createdAt,
		}

		rawDataLogs = append(rawDataLogs, rawLog)
	}

	if len(rawDataLogs) > 0 {
		// Bulk insert into PostgreSQL
		err = s.models.RawDataLog.BulkInsert(rawDataLogs)
		if err != nil {
			helpers.LogError(err, "Failed to insert raw data logs to PostgreSQL")
			return // Log the error and exit if bulk insert fails
		}
		helpers.LogInfo("Successfully inserted %d raw data logs into PostgreSQL", len(rawDataLogs))
	}

}
