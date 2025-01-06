package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/google/uuid"
)

// SigfoxKeepaliveLog represents a single Sigfox keepalive log entry.
type SigfoxKeepaliveLog struct {
	ID               int       `db:"id" json:"id"`
	RawID            uuid.UUID `db:"raw_id" json:"raw_id"`
	DeviceID         string    `db:"device_id" json:"device_id"`
	FirmwareVersion  float64   `db:"firmware_version" json:"firmware_version"`
	NetworkType      string    `db:"network_type" json:"network_type"`
	HappenedAt       time.Time `db:"happened_at" json:"happened_at"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	Timestamp        int64     `db:"timestamp" json:"timestamp"`
	IdleVoltage      int       `db:"idle_voltage" json:"idle_voltage"`
	Current          int       `db:"current" json:"current"`
	ResetCount       int       `db:"reset_count" json:"reset_count"`
	TemperatureMin   int       `db:"temperature_min" json:"temperature_min"`
	TemperatureMax   int       `db:"temperature_max" json:"temperature_max"`
	RadarError       int       `db:"radar_error" json:"radar_error"`
	TcveError        int       `db:"tcve_error" json:"tcve_error"`
	SettingsChecksum int       `db:"settings_checksum" json:"settings_checksum"`
}

// TableName returns the table name for the SigfoxKeepaliveLog model.
func (s *SigfoxKeepaliveLog) TableName() string {
	return "parking.sigfox_keepalive_logs"
}

// NewSigfoxKeepaliveLog constructs a SigfoxKeepaliveLog object from a provided map of data.
func NewSigfoxKeepaliveLog(pktData map[string]any) (*SigfoxKeepaliveLog, error) {
	// Parse raw_id as UUID
	rawUUIDStr, ok := pktData["raw_id"].(string)
	if !ok {
		return nil, helpers.WrapError(errors.New("invalid uuid format: expected string"))
	}

	rawUUID, err := uuid.Parse(rawUUIDStr)
	if err != nil {
		return nil, helpers.WrapError(fmt.Errorf("failed to parse uuid %s: %v", rawUUIDStr, err))
	}

	// Parse timestamp
	timestampFloat, ok := pktData["timestamp"].(float64)
	if !ok {
		return nil, helpers.WrapError(errors.New("invalid timestamp format: expected float64"))
	}
	timestamp := int64(timestampFloat)
	happenedAt := time.Unix(timestamp, 0).UTC()

	// Construct the SigfoxKeepaliveLog object
	log := &SigfoxKeepaliveLog{
		ID:               0, // ID is auto-incremented by the database.
		RawID:            rawUUID,
		DeviceID:         pktData["device_id"].(string),
		FirmwareVersion:  pktData["firmware_version"].(float64),
		NetworkType:      pktData["network_type"].(string),
		HappenedAt:       happenedAt,
		CreatedAt:        time.Now().UTC(),
		Timestamp:        timestamp,
		IdleVoltage:      int(pktData["idle_voltage"].(float64)),
		Current:          int(pktData["current"].(float64)),
		ResetCount:       int(pktData["reset_count"].(float64)),
		TemperatureMin:   int(pktData["temperature_min"].(float64)),
		TemperatureMax:   int(pktData["temperature_max"].(float64)),
		RadarError:       int(pktData["radar_error"].(float64)),
		TcveError:        int(pktData["tcve_error"].(float64)),
		SettingsChecksum: int(pktData["settings_checksum"].(float64)),
	}

	return log, nil
}

// BulkInsert inserts multiple SigfoxKeepaliveLog records in a single operation.
func (s *SigfoxKeepaliveLog) BulkInsert(keepaliveLogs []SigfoxKeepaliveLog) error {
	// Exit early if there are no records to insert
	if len(keepaliveLogs) == 0 {
		return nil
	}

	// Prepare slices for SQL values and arguments.
	values := make([]string, 0, len(keepaliveLogs))
	args := make([]interface{}, 0, len(keepaliveLogs)*15) // Adjust argument count based on the number of columns

	for i, log := range keepaliveLogs {
		// Create a placeholder for each record with indexed arguments
		values = append(values, fmt.Sprintf("( $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d )",
			i*15+1, i*15+2, i*15+3, i*15+4, i*15+5, i*15+6, i*15+7, i*15+8, i*15+9, i*15+10,
			i*15+11, i*15+12, i*15+13, i*15+14, i*15+15))

		// Append the actual values for each placeholder in the same order as the columns
		args = append(args,
			log.RawID, log.DeviceID, log.FirmwareVersion, log.NetworkType, log.HappenedAt, log.CreatedAt, log.Timestamp,
			log.IdleVoltage, log.Current, log.ResetCount, log.TemperatureMin, log.TemperatureMax, log.RadarError, log.TcveError,
			log.SettingsChecksum)
	}

	// Construct the SQL statement by joining the placeholders for each record
	query := fmt.Sprintf(
		`INSERT INTO %s (raw_id, device_id, firmware_version, network_type, happened_at, created_at, timestamp, idle_voltage, 
		current, reset_count, temperature_min, temperature_max, radar_error, tcve_error, settings_checksum) 
		VALUES %s
		`, s.TableName(), strings.Join(values, ","))

	// Execute the constructed query with the arguments
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk insert for Sigfox keepalive logs: %w", err)
	}

	return nil
}
