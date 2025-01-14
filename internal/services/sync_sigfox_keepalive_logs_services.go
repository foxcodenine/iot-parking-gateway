package services

import (
	"sort"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

func (s *Service) SyncSigfoxKeepaliveLogs() {
	// Retrieve keepalive log data from Redis and delete the key.

	items, err := s.cache.LRangeAndDelete("logs:sigfox-keepalive-logs")
	if err != nil {
		// Log error if Redis operations fail.
		helpers.LogError(err, "Error retrieving logs:sigfox-keepalive-logs from Redis")
		return
	}

	// Prepare a slice to hold the converted activity log entries.
	sigfoxKeepaliveLogs := make([]models.SigfoxKeepaliveLog, 0, len(items))

	// Iterate through each item retrieved from Redis.
	for _, item := range items {

		// Attempt to assert the type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			helpers.LogError(nil, "Invalid item type: expected map[string]any")
			continue
		}

		// Convert the map to an KeepaliveLog struct.
		keepaliveLog, err := models.NewSigfoxKeepaliveLog(itemMap)
		if err != nil {
			helpers.LogError(err, "")
		}

		// Append the successfully created activity log to the slice.
		sigfoxKeepaliveLogs = append(sigfoxKeepaliveLogs, *keepaliveLog)

	}

	// Sort keepalive logs by the HappenedAt field.
	sort.Slice(sigfoxKeepaliveLogs, func(i, j int) bool {
		return sigfoxKeepaliveLogs[i].HappenedAt.Before(sigfoxKeepaliveLogs[j].HappenedAt)
	})

	if len(sigfoxKeepaliveLogs) > 0 {

		err = s.models.SigfoxKeepaliveLog.BulkInsert(sigfoxKeepaliveLogs)

		// Log successful insertion of keepalive logs.
		if err != nil {
			helpers.LogError(err, "Failed to insert keepalive logs to PostgreSQL")
			return
		}
		helpers.LogInfo("Successfully inserted %d keepalive logs into PostgreSQL", len(sigfoxKeepaliveLogs))
	}
}
