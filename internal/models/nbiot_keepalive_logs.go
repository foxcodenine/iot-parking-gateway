package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/google/uuid"
)

// NbiotKeepaliveLog represents a single NB-IoT keepalive log entry.
type NbiotKeepaliveLog struct {
	ID                      int       `db:"id" json:"id"`                                                   // Auto-incrementing primary key
	RawID                   uuid.UUID `db:"raw_id" json:"raw_id"`                                           // ID linking to raw data source
	DeviceID                string    `db:"device_id" json:"device_id"`                                     // Device identifier, can be IMEI or UUID
	FirmwareVersion         float64   `db:"firmware_version" json:"firmware_version"`                       // Firmware version of the device
	NetworkType             string    `db:"network_type" json:"network_type"`                               // Network type (e.g., NB-IoT)
	HappenedAt              time.Time `db:"happened_at" json:"happened_at"`                                 // Time when the keepalive event occurred
	CreatedAt               time.Time `db:"created_at" json:"created_at"`                                   // Time when the record was created
	Timestamp               int64     `db:"timestamp" json:"timestamp"`                                     // Event timestamp in UNIX format
	IdleVoltage             int       `db:"idle_voltage" json:"idle_voltage"`                               // Idle voltage in V
	BatteryPercentage       int       `db:"battery_percentage" json:"battery_percentage"`                   // Battery percentage
	Current                 int       `db:"current" json:"current"`                                         // Current in mA
	ResetCount              int       `db:"reset_count" json:"reset_count"`                                 // Reset count
	ManualCalibration       bool      `db:"manual_calibration" json:"manual_calibration"`                   // Manual calibration status
	TemperatureMin          int       `db:"temperature_min" json:"temperature_min"`                         // Minimum temperature
	TemperatureMax          int       `db:"temperature_max" json:"temperature_max"`                         // Maximum temperature
	RadarError              int       `db:"radar_error" json:"radar_error"`                                 // Radar error count
	MagError                int       `db:"mag_error" json:"mag_error"`                                     // Magnetometer error count
	TcveError               int       `db:"tcve_error" json:"tcve_error"`                                   // TCVE error count
	BleSecurityIssues       int       `db:"ble_security_issues" json:"ble_security_issues"`                 // BLE security issues count
	RadarCumulativeTotal    int       `db:"radar_cumulative_total" json:"radar_cumulative_total"`           // Radar cumulative total
	MagTotal                int       `db:"mag_total" json:"mag_total"`                                     // Magnetometer total value
	NetworkRegistrationOk   int       `db:"network_registration_ok" json:"network_registration_ok"`         // Registration successful count
	NetworkRegistrationNok  int       `db:"network_registration_nok" json:"network_registration_nok"`       // Registration failed count
	RssiAverage             int       `db:"rssi_average" json:"rssi_average"`                               // RSSI average value
	NetworkMessageAttempts  int       `db:"network_message_attempts" json:"network_message_attempts"`       // Message attempts count
	NetworkAck1ds           int       `db:"network_ack_1ds" json:"network_ack_1ds"`                         // Network ACK 1DS count
	Network1ackDs           int       `db:"network_1ack_ds" json:"network_1ack_ds"`                         // Network 1ACK DS count
	Network1ack1ds          int       `db:"network_1ack_1ds" json:"network_1ack_1ds"`                       // Network 1ACK 1DS count
	TcvrDeepSleepMin        int       `db:"tcvr_deep_sleep_min" json:"tcvr_deep_sleep_min"`                 // Deep sleep min time
	TcvrDeepSleepMax        int       `db:"tcvr_deep_sleep_max" json:"tcvr_deep_sleep_max"`                 // Deep sleep max time
	TcvrDeepSleepAverage    int       `db:"tcvr_deep_sleep_average" json:"tcvr_deep_sleep_average"`         // Deep sleep average time
	SettingsChecksum        int       `db:"settings_checksum" json:"settings_checksum"`                     // Settings checksum value
	SocketError             int       `db:"socket_error" json:"socket_error"`                               // Socket error count
	T3324                   int       `db:"t3324" json:"t3324"`                                             // Timer T3324 value
	T3412                   int       `db:"t3412" json:"t3412"`                                             // Timer T3412 value
	TimeSyncRandByte        int       `db:"time_sync_rand_byte" json:"time_sync_rand_byte"`                 // Time sync random byte
	TimeSyncCurrentUnixTime int64     `db:"time_sync_current_unix_time" json:"time_sync_current_unix_time"` // Time sync current UNIX time
}

// TableName returns the table name for the NbiotKeepaliveLog model.
func (n *NbiotKeepaliveLog) TableName() string {
	return "parking.nbiot_keepalive_logs"
}

// NewNbiotKeepaliveLog constructs an NbiotKeepaliveLog object from a provided map of data.
// It handles data type conversions and populates the fields accordingly.
func NewNbiotKeepaliveLog(pktData map[string]any) (*NbiotKeepaliveLog, error) {
	// Parse the raw UUID field from a string to uuid.UUID.
	rawUUIDStr, ok := pktData["raw_id"].(string)
	if !ok {
		return nil, helpers.WrapError(errors.New("invalid uuid format: expected string"))
	}

	rawUUID, err := uuid.Parse(rawUUIDStr)
	if err != nil {
		return nil, helpers.WrapError(fmt.Errorf("failed to parse uuid %s: %v", rawUUIDStr, err))
	}

	// Parse the timestamp, expected as an int64.

	timestampFloat, ok := pktData["timestamp"].(float64)
	if !ok {
		return nil, helpers.WrapError(errors.New("invalid timestamp format: expected float64"))
	}
	// Convert float64 timestamp to int64 and then to time.Time.
	timestamp := int64(timestampFloat)
	// Convert int64 timestamp to time.Time.
	happenedAt := time.Unix(timestamp, 0).UTC()

	// Construct and return the NbiotKeepaliveLog object with the parsed and converted data.
	nbiotKeepaliveLog := &NbiotKeepaliveLog{
		ID:                     0, // ID is auto-incremented by the database.
		RawID:                  rawUUID,
		DeviceID:               pktData["device_id"].(string),
		FirmwareVersion:        pktData["firmware_version"].(float64),
		NetworkType:            pktData["network_type"].(string),
		HappenedAt:             happenedAt,
		CreatedAt:              time.Now().UTC(), // Default to the current time in UTC.
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
		SocketError:            int(pktData["socket_error"].(float64)),
		T3324:                  int(pktData["t3324"].(float64)),
		T3412:                  int(pktData["t3412"].(float64)),
		// TimeSyncRandByte:       int(pktData["time_sync_rand_byte"].(float64)),
		// 	// TimeSyncCurrentUnixTime: int64(pktData["time_sync_current_unix_time"].(float64)), // New field.
	}

	if pktData["time_sync_rand_byte"] != nil {
		nbiotKeepaliveLog.TimeSyncRandByte = int(pktData["time_sync_rand_byte"].(float64))
	}
	if pktData["time_sync_current_unix_time"] != nil {
		nbiotKeepaliveLog.TimeSyncCurrentUnixTime = int64(pktData["time_sync_current_unix_time"].(float64))
	}

	return nbiotKeepaliveLog, nil
}

// BulkInsert inserts multiple NbiotKeepaliveLog records in a single operation.
func (n *NbiotKeepaliveLog) BulkInsert(keepaliveLogs []NbiotKeepaliveLog) error {
	// Exit early if there are no records to insert
	if len(keepaliveLogs) == 0 {
		return nil
	}

	// Prepare slices for SQL values and arguments.
	values := make([]string, 0, len(keepaliveLogs))
	args := make([]interface{}, 0, len(keepaliveLogs)*36) // Adjust the argument count based on the number of columns

	for i, log := range keepaliveLogs {
		// Create a placeholder for each record with indexed arguments
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*36+1, i*36+2, i*36+3, i*36+4, i*36+5, i*36+6, i*36+7, i*36+8, i*36+9, i*36+10,
			i*36+11, i*36+12, i*36+13, i*36+14, i*36+15, i*36+16, i*36+17, i*36+18, i*36+19, i*36+20,
			i*36+21, i*36+22, i*36+23, i*36+24, i*36+25, i*36+26, i*36+27, i*36+28, i*36+29, i*36+30,
			i*36+31, i*36+32, i*36+33, i*36+34, i*36+35, i*36+36))

		// Append the actual values for each placeholder in the same order as the columns
		args = append(args,
			log.RawID, log.DeviceID, log.FirmwareVersion, log.NetworkType, log.HappenedAt, log.CreatedAt, log.Timestamp,
			log.IdleVoltage, log.BatteryPercentage, log.Current, log.ResetCount, log.ManualCalibration, log.TemperatureMin,
			log.TemperatureMax, log.RadarError, log.MagError, log.TcveError, log.BleSecurityIssues, log.RadarCumulativeTotal,
			log.MagTotal, log.NetworkRegistrationOk, log.NetworkRegistrationNok, log.RssiAverage, log.NetworkMessageAttempts,
			log.NetworkAck1ds, log.Network1ackDs, log.Network1ack1ds, log.TcvrDeepSleepMin, log.TcvrDeepSleepMax,
			log.TcvrDeepSleepAverage, log.SettingsChecksum, log.SocketError, log.T3324, log.T3412, log.TimeSyncRandByte,
			log.TimeSyncCurrentUnixTime)
	}

	// Construct the SQL statement by joining the placeholders for each record
	query := fmt.Sprintf(
		`INSERT INTO %s (raw_id, device_id, firmware_version, network_type, happened_at, created_at, timestamp,
		idle_voltage, battery_percentage, current, reset_count, manual_calibration, temperature_min, temperature_max,
		radar_error, mag_error, tcve_error, ble_security_issues, radar_cumulative_total, mag_total,
		network_registration_ok, network_registration_nok, rssi_average, network_message_attempts, network_ack_1ds,
		network_1ack_ds, network_1ack_1ds, tcvr_deep_sleep_min, tcvr_deep_sleep_max, tcvr_deep_sleep_average,
		settings_checksum, socket_error, t3324, t3412, time_sync_rand_byte, time_sync_current_unix_time) VALUES %s`,
		n.TableName(), strings.Join(values, ", "))

	// Execute the constructed query with the arguments
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk insert for keepalive logs: %w", err)
	}

	return nil
}
