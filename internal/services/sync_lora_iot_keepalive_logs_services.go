package services

import (
	"sort"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

func (s *Service) SyncLoraKeepaliveLogs() {
	// Retrieve keepalive log data from Redis and delete the key.

	items, err := s.cache.LRangeAndDelete("logs:lora-keepalive-logs")
	if err != nil {
		// Log error if Redis operations fail.
		s.errorLog.Printf("Error retrieving logs:lora-keepalive-logs from Redis: %v", err)
		return
	}

	// Prepare a slice to hold the converted activity log entries.
	loraKeepaliveLogs := make([]models.LoraKeepaliveLog, 0, len(items))

	// Iterate through each item retrieved from Redis.
	for _, item := range items {

		// Attempt to assert the type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			s.errorLog.Println("Invalid item type: expected map[string]any")
			continue
		}

		// Convert the map to an KeepaliveLog struct.
		keepaliveLog, err := models.NewLoraKeepaliveLog(itemMap)
		if err != nil {
			helpers.LogError(err, "")
		}

		// Append the successfully created activity log to the slice.
		loraKeepaliveLogs = append(loraKeepaliveLogs, *keepaliveLog)

	}

	// Sort keepalive logs by the HappenedAt field.
	sort.Slice(loraKeepaliveLogs, func(i, j int) bool {
		return loraKeepaliveLogs[i].HappenedAt.Before(loraKeepaliveLogs[j].HappenedAt)
	})

	if len(loraKeepaliveLogs) > 0 {

		err = s.models.LoraKeepaliveLog.BulkInsert(loraKeepaliveLogs)

		// Log successful insertion of keepalive logs.
		if err != nil {
			s.errorLog.Printf("Failed to insert keepalive logs to PostgreSQL: %v", err)
			return
		}
		s.infoLog.Printf("Successfully inserted %d keepalive logs into PostgreSQL", len(loraKeepaliveLogs))
	}
}