package services

import (
	"sort"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

func (s *Service) SyncNBIoTKeepaliveLogs() {
	// Retrieve keepalive log data from Redis and delete the key.

	items, err := s.cache.LRangeAndDelete("logs:nb-keepalive-logs")
	if err != nil {
		// Log error if Redis operations fail.
		helpers.LogError(err, "Error retrieving logs:nb-keepalive-logs from Redis")
		return
	}

	// Prepare a slice to hold the converted activity log entries.
	nbIotKeepaliveLogs := make([]models.NbiotKeepaliveLog, 0, len(items))

	// Iterate through each item retrieved from Redis.
	for _, item := range items {

		// Attempt to assert the type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			helpers.LogError(nil, "Invalid item type: expected map[string]any")
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

	// Sort keepalive logs by the HappenedAt field.
	sort.Slice(nbIotKeepaliveLogs, func(i, j int) bool {
		return nbIotKeepaliveLogs[i].HappenedAt.Before(nbIotKeepaliveLogs[j].HappenedAt)
	})

	if len(nbIotKeepaliveLogs) > 0 {

		err = s.models.NbiotKeepaliveLog.BulkInsert(nbIotKeepaliveLogs)

		// Log successful insertion of keepalive logs.
		if err != nil {
			helpers.LogError(err, "Failed to insert keepalive logs to PostgreSQL")
			return
		}
		helpers.LogInfo("Successfully inserted %d keepalive logs into PostgreSQL", len(nbIotKeepaliveLogs))
	}
}
