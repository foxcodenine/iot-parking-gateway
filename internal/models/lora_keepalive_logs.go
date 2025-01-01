package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/google/uuid"
)

// LoraKeepaliveLog represents a single LoRa keepalive log entry.
type LoraKeepaliveLog struct {
	ID                      int       `db:"id" json:"id"`
	RawID                   uuid.UUID `db:"raw_id" json:"raw_id"`
	DeviceID                string    `db:"device_id" json:"device_id"`
	FirmwareVersion         float64   `db:"firmware_version" json:"firmware_version"`
	NetworkType             string    `db:"network_type" json:"network_type"`
	HappenedAt              time.Time `db:"happened_at" json:"happened_at"`
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
	Timestamp               int64     `db:"timestamp" json:"timestamp"`
	IdleVoltage             int       `db:"idle_voltage" json:"idle_voltage"`
	BatteryPercentage       int       `db:"battery_percentage" json:"battery_percentage"`
	Current                 int       `db:"current" json:"current"`
	ResetCount              int       `db:"reset_count" json:"reset_count"`
	ManualCalibration       bool      `db:"manual_calibration" json:"manual_calibration"`
	TemperatureMin          int       `db:"temperature_min" json:"temperature_min"`
	TemperatureMax          int       `db:"temperature_max" json:"temperature_max"`
	RadarError              int       `db:"radar_error" json:"radar_error"`
	MagError                int       `db:"mag_error" json:"mag_error"`
	TcveError               int       `db:"tcve_error" json:"tcve_error"`
	BleSecurityIssues       int       `db:"ble_security_issues" json:"ble_security_issues"`
	RadarCumulativeTotal    int       `db:"radar_cumulative_total" json:"radar_cumulative_total"`
	MagTotal                int       `db:"mag_total" json:"mag_total"`
	NetworkRegistrationOk   int       `db:"network_registration_ok" json:"network_registration_ok"`
	NetworkRegistrationNok  int       `db:"network_registration_nok" json:"network_registration_nok"`
	RssiAverage             int       `db:"rssi_average" json:"rssi_average"`
	NetworkMessageAttempts  int       `db:"network_message_attempts" json:"network_message_attempts"`
	NetworkAck1ds           int       `db:"network_ack_1ds" json:"network_ack_1ds"`
	Network1ackDs           int       `db:"network_1ack_ds" json:"network_1ack_ds"`
	Network1ack1ds          int       `db:"network_1ack_1ds" json:"network_1ack_1ds"`
	TcvrDeepSleepMin        int       `db:"tcvr_deep_sleep_min" json:"tcvr_deep_sleep_min"`
	TcvrDeepSleepMax        int       `db:"tcvr_deep_sleep_max" json:"tcvr_deep_sleep_max"`
	TcvrDeepSleepAverage    int       `db:"tcvr_deep_sleep_average" json:"tcvr_deep_sleep_average"`
	SettingsChecksum        int       `db:"settings_checksum" json:"settings_checksum"`
	TimeSyncRandByte        int       `db:"time_sync_rand_byte" json:"time_sync_rand_byte"`
	TimeSyncCurrentUnixTime *int64    `db:"time_sync_current_unix_time" json:"time_sync_current_unix_time"` // Nullable field
}

// TableName returns the table name for the LoraKeepaliveLog model.
func (l *LoraKeepaliveLog) TableName() string {
	return "parking.lora_keepalive_logs"
}

// NewLoraKeepaliveLog constructs a LoraKeepaliveLog object from a provided map of data.
func NewLoraKeepaliveLog(pktData map[string]any) (*LoraKeepaliveLog, error) {
	rawUUIDStr, ok := pktData["raw_id"].(string)
	if !ok {
		return nil, helpers.WrapError(errors.New("invalid uuid format: expected string"))
	}

	rawUUID, err := uuid.Parse(rawUUIDStr)
	if err != nil {
		return nil, helpers.WrapError(fmt.Errorf("failed to parse uuid %s: %v", rawUUIDStr, err))
	}

	timestampFloat, ok := pktData["timestamp"].(float64)
	if !ok {
		return nil, helpers.WrapError(errors.New("invalid timestamp format: expected float64"))
	}
	timestamp := int64(timestampFloat)
	happenedAt := time.Unix(timestamp, 0).UTC()

	log := &LoraKeepaliveLog{
		ID:                     0, // ID is auto-incremented by the database.
		RawID:                  rawUUID,
		DeviceID:               pktData["device_id"].(string),
		FirmwareVersion:        pktData["firmware_version"].(float64),
		NetworkType:            pktData["network_type"].(string),
		HappenedAt:             happenedAt,
		CreatedAt:              time.Now().UTC(),
		Timestamp:              timestamp,
		IdleVoltage:            int(pktData["idle_voltage"].(float64)),
		BatteryPercentage:      int(pktData["battery_percentage"].(float64)),
		Current:                int(pktData["current"].(float64)),
		ResetCount:             int(pktData["reset_count"].(float64)),
		ManualCalibration:      pktData["manual_calibration"].(float64) != 0,
		TemperatureMin:         int(pktData["temperature_min"].(float64)),
		TemperatureMax:         int(pktData["temperature_max"].(float64)),
		RadarError:             int(pktData["radar_error"].(float64)),
		MagError:               int(pktData["mag_error"].(float64)),
		TcveError:              int(pktData["tcve_error"].(float64)),
		BleSecurityIssues:      int(pktData["ble_security_issues"].(float64)),
		RadarCumulativeTotal:   int(pktData["radar_cumulative_total"].(float64)),
		MagTotal:               int(pktData["mag_total"].(float64)),
		NetworkRegistrationOk:  int(pktData["network_registration_ok"].(float64)),
		NetworkRegistrationNok: int(pktData["network_registration_nok"].(float64)),
		RssiAverage:            int(pktData["rssi_average"].(float64)),
		NetworkMessageAttempts: int(pktData["network_message_attempts"].(float64)),
		NetworkAck1ds:          int(pktData["network_ack_1ds"].(float64)),
		Network1ackDs:          int(pktData["network_1ack_ds"].(float64)),
		Network1ack1ds:         int(pktData["network_1ack_1ds"].(float64)),
		TcvrDeepSleepMin:       int(pktData["tcvr_deep_sleep_min"].(float64)),
		TcvrDeepSleepMax:       int(pktData["tcvr_deep_sleep_max"].(float64)),
		TcvrDeepSleepAverage:   int(pktData["tcvr_deep_sleep_average"].(float64)),
		SettingsChecksum:       int(pktData["settings_checksum"].(float64)),
		TimeSyncRandByte:       int(pktData["time_sync_rand_byte"].(float64)),
	}

	// Optionally handle TimeSyncCurrentUnixTime if it's present
	if ts, exists := pktData["time_sync_current_unix_time"]; exists {
		if tsFloat, ok := ts.(float64); ok {
			tsInt64 := int64(tsFloat)
			log.TimeSyncCurrentUnixTime = &tsInt64
		} else {
			return nil, helpers.WrapError(errors.New("invalid time_sync_current_unix_time format: expected float64"))
		}
	}

	return log, nil
}

// BulkInsert inserts multiple LoraKeepaliveLog records in a single operation.
func (l *LoraKeepaliveLog) BulkInsert(keepaliveLogs []LoraKeepaliveLog) error {
	// Exit early if there are no records to insert
	if len(keepaliveLogs) == 0 {
		return nil
	}

	// Prepare slices for SQL values and arguments.
	values := make([]string, 0, len(keepaliveLogs))
	args := make([]interface{}, 0, len(keepaliveLogs)*36) // Adjust the argument count based on the number of columns

	for i, log := range keepaliveLogs {
		// Create a placeholder for each record with indexed arguments
		values = append(values, fmt.Sprintf("( $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*33+1, i*33+2, i*33+3, i*33+4, i*33+5, i*33+6, i*33+7, i*33+8, i*33+9, i*33+10,
			i*33+11, i*33+12, i*33+13, i*33+14, i*33+15, i*33+16, i*33+17, i*33+18, i*33+19, i*33+20,
			i*33+21, i*33+22, i*33+23, i*33+24, i*33+25, i*33+26, i*33+27, i*33+28, i*33+29, i*33+30,
			i*33+31, i*33+32, i*33+33))

		// Append the actual values for each placeholder in the same order as the columns
		args = append(args,
			log.RawID, log.DeviceID, log.FirmwareVersion, log.NetworkType, log.HappenedAt, log.CreatedAt, log.Timestamp,
			log.IdleVoltage, log.BatteryPercentage, log.Current, log.ResetCount, log.ManualCalibration, log.TemperatureMin,
			log.TemperatureMax, log.RadarError, log.MagError, log.TcveError, log.BleSecurityIssues, log.RadarCumulativeTotal,
			log.MagTotal, log.NetworkRegistrationOk, log.NetworkRegistrationNok, log.RssiAverage, log.NetworkMessageAttempts,
			log.NetworkAck1ds, log.Network1ackDs, log.Network1ack1ds, log.TcvrDeepSleepMin, log.TcvrDeepSleepMax,
			log.TcvrDeepSleepAverage, log.SettingsChecksum, log.TimeSyncRandByte, log.TimeSyncCurrentUnixTime)
	}

	// Construct the SQL statement by joining the placeholders for each record
	query := fmt.Sprintf(
		`INSERT INTO %s (raw_id, device_id, firmware_version, network_type, happened_at, created_at, timestamp, idle_voltage, 
		battery_percentage, current, reset_count, manual_calibration, temperature_min, temperature_max, radar_error, mag_error, 
		tcve_error, ble_security_issues, radar_cumulative_total, mag_total, network_registration_ok, network_registration_nok, 
		rssi_average, network_message_attempts, network_ack_1ds, network_1ack_ds, network_1ack_1ds, tcvr_deep_sleep_min, 
		tcvr_deep_sleep_max, tcvr_deep_sleep_average, settings_checksum, time_sync_rand_byte, time_sync_current_unix_time) 
		VALUES %s
		`, l.TableName(), strings.Join(values, ","))

	// Execute the constructed query with the arguments
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk insert for keepalive logs: %w", err)
	}

	return nil
}
