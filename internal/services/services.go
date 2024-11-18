package services

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/google/uuid"
)

type Service struct {
	models   models.Models
	cache    *cache.RedisCache
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewService(m models.Models, rc *cache.RedisCache, il, el *log.Logger) *Service {
	return &Service{
		models:   m,
		cache:    rc,
		infoLog:  il,
		errorLog: el,
	}
}

// SyncRawLogs processes and synchronizes raw logs from Redis to PostgreSQL.
func (s *Service) SyncRawLogs() {

	items, err := s.cache.LRangeAndDelete("raw-data-logs")
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

	// Bulk insert into PostgreSQL
	err = s.models.RawDataLog.BulkInsert(rawDataLogs)
	if err != nil {
		s.errorLog.Printf("Failed to insert raw data logs to PostgreSQL: %v", err)
		return // Log the error and exit if bulk insert fails
	}

	s.infoLog.Printf("Successfully inserted %d raw data logs into PostgreSQL", len(rawDataLogs))

}

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

	// Attempt to bulk insert all activity logs into PostgreSQL.
	if err = s.models.ActivityLog.BulkInsert(activityLogs); err != nil {
		s.errorLog.Printf("Failed to insert activity logs to PostgreSQL: %v", err)
		return
	}

	// Log successful insertion of activity logs.
	s.infoLog.Printf("Successfully inserted %d activity logs into PostgreSQL", len(activityLogs))

	// If activity logs are successfully inserted, proceed to update devices.
	if err = s.models.Device.BulkUpdateDevices(activityLogs); err != nil {
		s.errorLog.Printf("Failed to update device records: %v", err)
		return
	}

	// Log successful update of devices.
	s.infoLog.Printf("Successfully updated device records based on activity logs")
}

func (s *Service) SyncNBIoTKeepaliveLogs() {
	// Retrieve keepalive log data from Redis and delete the key.
	items, err := s.cache.LRangeAndDelete("nb-keepalive-logs")
	if err != nil {
		// Log error if Redis operations fail.
		s.errorLog.Printf("Error retrieving items from Redis: %v", err)
		return
	}

	// Prepare a slice to hold the converted activity log entries.
	nbIotKeepaliveLogs := make([]models.NbiotKeepaliveLog, 0, len(items))

	// Iterate through each item retrieved from Redis.
	for _, item := range items {

		// Attempt to assert the type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			s.errorLog.Println("Invalid item type: expected map[string]any")
			continue
		}

		// Convert the map to an KeepaliveLog struct.
		keepaliveLog, err := models.NewNbiotKeepaliveLog(itemMap)
		if err != nil {
			helpers.LogError(err, "")
		}

		// Append the successfully created activity log to the slice.
		nbIotKeepaliveLogs = append(nbIotKeepaliveLogs, *keepaliveLog)

	}

	s.models.NbiotKeepaliveLog.BulkInsert(nbIotKeepaliveLogs)
}

// RegisterNewDevices retrieves device IDs from Redis, creates new device records in the database, and logs the outcome.
func (s *Service) RegisterNewDevices() {
	// Retrieve device IDs from Redis and delete the set afterwards.
	deviceEntries, err := s.cache.SMembersDel("device-to-create")
	if err != nil {
		s.errorLog.Printf("Failed to retrieve device IDs from Redis: %v", err)
		return
	}

	// Process each device entry to create new device records.
	for _, entry := range deviceEntries {
		// Convert interface to string, ensuring it represents a device ID correctly.
		deviceInfo, ok := entry.(string)
		if !ok {
			s.errorLog.Printf("Invalid device ID format: expected string but got %T", entry)
			continue
		}

		// Extract network type and device ID from the string.
		spaceIndex := strings.Index(deviceInfo, " ")
		if spaceIndex == -1 {
			s.errorLog.Printf("Device information string does not contain a space: %s", deviceInfo)
			continue
		}

		networkType := deviceInfo[:spaceIndex]
		deviceID := deviceInfo[spaceIndex+1:]

		// Define a new device model instance.
		newDevice := models.Device{
			DeviceID:    deviceID,
			Name:        "new " + deviceID,
			NetworkType: networkType,
		}

		// Attempt to create a new device record in the database.
		_, err = s.models.Device.Create(&newDevice)
		if err != nil {
			s.errorLog.Printf("Failed to create a new device record for ID %s: %v", deviceID, err)
			continue // Proceed to the next ID instead of stopping.
		}

	}
	if len(deviceEntries) == 1 {
		s.infoLog.Println("Successfully added 1 device to PostgreSQL")
	} else {
		s.infoLog.Printf("Successfully added %d devices to PostgreSQL", len(deviceEntries))
	}
}

// PopulateDeviceBloomFilter retrieves all devices from the database and populates a Bloom Filter with their network type and device IDs.
func (s *Service) PopulateDeviceBloomFilter() {
	// Fetch all devices from the database.

	tmp := models.Device{}
	devices, err := tmp.GetAll()
	if err != nil {
		s.errorLog.Printf("Error retrieving devices from Postgres: %v", err)
		return
	}

	// Initialize a slice to store composite keys for the Bloom Filter.
	var keyStrings []string

	// Iterate over each device to create a composite key of network type and device ID.
	for _, d := range devices {
		key := fmt.Sprintf("%s %s", d.NetworkType, d.DeviceID)
		keyStrings = append(keyStrings, key)
	}

	// Add all composite keys to the Bloom Filter named "device-id".
	s.cache.AddItemsToBloomFilter("device-id", keyStrings)
}
