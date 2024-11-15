package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RawDataLog struct {
	ID              uuid.UUID `db:"id" json:"id"`
	DeviceID        string    `db:"device_id" json:"device_id"`
	FirmwareVersion float64   `db:"firmware_version" json:"firmware_version"`
	NetworkType     string    `db:"network_type" json:"network_type"`
	RawData         string    `db:"raw_data" json:"raw_data"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

func (r *RawDataLog) TableName() string {
	return "parking.raw_data_logs"
}

func (r *RawDataLog) BulkInsert(rawDataLogs []RawDataLog) error {
	if len(rawDataLogs) == 0 {
		return nil
	}

	// Prepare slices for SQL values and arguments.
	values := make([]string, 0, len(rawDataLogs))      // Holds each row as a set of placeholders
	args := make([]interface{}, 0, len(rawDataLogs)*5) // Holds each actual value to be inserted

	for i, log := range rawDataLogs {
		// Create a placeholder for each record, referencing argument indices for this row.
		// The format "($1, $2, ..., $5)" is modified per row to match the argument position.
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))

		// Add each field value to the args slice in the same order as the placeholders.
		args = append(args, log.ID, log.DeviceID, log.FirmwareVersion, log.NetworkType, log.RawData)
	}

	// Construct the full SQL query string for bulk insertion.
	query := fmt.Sprintf("INSERT INTO %s (id, device_id, firmware_version, network_type, raw_data) VALUES %s",
		r.TableName(), strings.Join(values, ", "))

	// Execute the bulk insert query
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk insert: %w", err)
	}

	return nil
}
