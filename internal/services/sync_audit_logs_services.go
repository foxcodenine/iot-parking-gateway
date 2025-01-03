package services

import (
	"sort"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

func (s *Service) SyncAuditLogs() {

	items, err := s.cache.LRangeAndDelete("logs:audit-logs")

	if err != nil {
		// Log error if Redis operations fail.
		s.errorLog.Printf("Error retrieving logs:audit-logs from Redis: %v", err)
		return
	}

	auditLogs := make([]models.AuditLog, 0, len(items))

	for _, item := range items {

		// Attempt to assert the type to a map[string]any (JSON-like structure).
		itemMap, ok := item.(map[string]any)
		if !ok {
			s.errorLog.Println("Invalid item type: expected map[string]any")
			continue
		}

		// Convert the map to an KeepaliveLog struct.
		auditLog, err := models.NewAuditLog(itemMap)
		if err != nil {
			helpers.LogError(err, "")
		}

		// Append the successfully created activity log to the slice.
		auditLogs = append(auditLogs, *auditLog)

	}

	// Sort keepalive logs by the HappenedAt field.
	sort.Slice(auditLogs, func(i, j int) bool {
		return auditLogs[i].HappenedAt.Before(auditLogs[j].HappenedAt)
	})

	if len(auditLogs) > 0 {

		err = s.models.AuditLog.BulkInsert(auditLogs)

		// Log successful insertion of keepalive logs.
		if err != nil {
			s.errorLog.Printf("Failed to insert auditLogs logs to PostgreSQL: %v", err)
			return
		}
		s.infoLog.Printf("Successfully inserted %d auditLogs logs into PostgreSQL", len(auditLogs))
	}
}
