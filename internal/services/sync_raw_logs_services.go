package services

import (
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/google/uuid"
)

// SyncRawLogs processes and synchronizes raw logs from Redis to PostgreSQL.
func (s *Service) SyncRawLogs() {

	items, err := s.cache.LRangeAndDelete("logs:raw-data-logs")
	if err != nil {
		s.errorLog.Printf("Error retrieving items from Redis: %v", err)
		return
	}

	var rawDataLogs []models.RawDataLog

	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			s.errorLog.Println("Invalid item type: expected map[string]interface{}")
			continue
		}

		// Convert id field from string to uuid.UUID
		uuidStr, ok := itemMap["id"].(string)
		if !ok {
			s.errorLog.Println("Invalid UUID format: expected string")
			continue
		}

		uuidValue, err := uuid.Parse(uuidStr)
		if err != nil {
			s.errorLog.Printf("Failed to parse UUID %s: %v\n", uuidStr, err)
			continue
		}

		// Parse created_at field from string to time.Time
		createAtStr, ok := itemMap["created_at"].(string)
		if !ok {
			s.errorLog.Println("Invalid created_at format: expected string")
			continue
		}

		createdAt, err := time.Parse(time.RFC3339, createAtStr)
		if err != nil {
			s.errorLog.Printf("Failed to parse created_at %s: %v\n", createAtStr, err)
			continue
		}

		// Check for and convert other fields safely
		deviceID, ok := itemMap["device_id"].(string)
		if !ok {
			s.errorLog.Println("Invalid device_id format: expected string")
			continue
		}

		firmwareVersion, ok := itemMap["firmware_version"].(float64) // Redis stores numbers as float64
		if !ok {
			s.errorLog.Println("Invalid firmware_version format: expected float64")
			continue
		}

		networkType, ok := itemMap["network_type"].(string)
		if !ok {
			s.errorLog.Println("Invalid network_type format: expected string")
			continue
		}

		rawData, ok := itemMap["raw_data"].(string)
		if !ok {
			s.errorLog.Println("Invalid raw_data format: expected string")
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
			s.errorLog.Printf("Failed to insert raw data logs to PostgreSQL: %v", err)
			return // Log the error and exit if bulk insert fails
		}
		s.infoLog.Printf("Successfully inserted %d raw data logs into PostgreSQL", len(rawDataLogs))
	}

}
