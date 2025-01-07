package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SigfoxSettingLog represents a single Sigfox setting log entry.
type SigfoxSettingLog struct {
	ID                                         int       `db:"id" json:"id"`
	RawID                                      uuid.UUID `db:"raw_id" json:"raw_id"`
	DeviceID                                   string    `db:"device_id" json:"device_id"`
	FirmwareVersion                            float64   `db:"firmware_version" json:"firmware_version"`
	NetworkType                                string    `db:"network_type" json:"network_type"`
	HappenedAt                                 time.Time `db:"happened_at" json:"happened_at"`
	CreatedAt                                  time.Time `db:"created_at" json:"created_at"`
	Timestamp                                  int64     `db:"timestamp" json:"timestamp"`
	DeviceMode                                 int       `db:"device_mode" json:"device_mode"`
	DeviceEnable                               int       `db:"device_enable" json:"device_enable"`
	RadarCarCalLoTh                            int       `db:"radar_car_cal_lo_th" json:"radar_car_cal_lo_th"`
	RadarCarCalHiTh                            int       `db:"radar_car_cal_hi_th" json:"radar_car_cal_hi_th"`
	RadarCarDeltaTh                            int       `db:"radar_car_delta_th" json:"radar_car_delta_th"`
	DownlinkEn7BitsRepeatedOccupancyPeriodMins int       `db:"downlink_en_7_bits_repeated_occupancy_period_mins" json:"downlink_en_7_bits_repeated_occupancy_period_mins"`
}

// TableName returns the table name for the SigfoxSettingLog model.
func (s *SigfoxSettingLog) TableName() string {
	return "parking.sigfox_setting_logs"
}

// NewSigfoxSettingLog constructs a SigfoxSettingLog object from a provided map of data.
// It handles data type conversions and populates the fields accordingly.
func NewSigfoxSettingLog(pktData map[string]any) (*SigfoxSettingLog, error) {
	// Parse the raw UUID field from a string to uuid.UUID.
	rawUUIDStr, ok := pktData["raw_id"].(string)
	if !ok {
		return nil, errors.New("invalid uuid format: expected string")
	}

	rawUUID, err := uuid.Parse(rawUUIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid %s: %v", rawUUIDStr, err)
	}

	// Parse the timestamp, expected as a float64.
	timestampFloat, ok := pktData["timestamp"].(float64)
	if !ok {
		return nil, errors.New("invalid timestamp format: expected float64")
	}
	// Convert float64 timestamp to int64 and then to time.Time.
	timestamp := int64(timestampFloat)
	happenedAt := time.Unix(timestamp, 0).UTC()

	// Construct and return the SigfoxSettingLog object with the parsed and converted data.
	log := &SigfoxSettingLog{
		RawID:           rawUUID,
		DeviceID:        pktData["device_id"].(string),
		FirmwareVersion: pktData["firmware_version"].(float64),
		NetworkType:     pktData["network_type"].(string),
		HappenedAt:      happenedAt,
		CreatedAt:       time.Now().UTC(),
		Timestamp:       timestamp,
		DeviceMode:      int(pktData["device_mode"].(float64)),
		DeviceEnable:    int(pktData["device_enable"].(float64)),
		RadarCarCalLoTh: int(pktData["radar_car_cal_lo_th"].(float64)),
		RadarCarCalHiTh: int(pktData["radar_car_cal_hi_th"].(float64)),
		RadarCarDeltaTh: int(pktData["radar_car_delta_th"].(float64)),
		DownlinkEn7BitsRepeatedOccupancyPeriodMins: int(pktData["downlink_en_7_bits_repeated_occupancy_period_mins"].(float64)),
	}

	return log, nil
}

// BulkInsert inserts multiple SigfoxSettingLog records in a single operation.
func (s *SigfoxSettingLog) BulkInsert(settingLogs []SigfoxSettingLog) error {
	// Exit early if there are no records to insert
	if len(settingLogs) == 0 {
		return nil
	}

	// Determine the number of fields to be inserted for each log
	numFields := 13 // Adjust this based on the actual number of columns in your table

	// Prepare slices for SQL values and arguments.
	values := make([]string, 0, len(settingLogs))
	args := make([]interface{}, 0, len(settingLogs)*numFields) // Adjust the argument count based on the number of columns

	for i, log := range settingLogs {
		// Create a placeholder for each record with indexed arguments
		placeholders := make([]string, numFields)
		for j := range placeholders {
			placeholders[j] = fmt.Sprintf("$%d", i*numFields+j+1)
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))

		// Append the actual values for each placeholder in the same order as the columns
		args = append(args,
			log.RawID, log.DeviceID, log.FirmwareVersion, log.NetworkType, log.HappenedAt, log.CreatedAt, log.Timestamp,
			log.DeviceMode, log.DeviceEnable, log.RadarCarCalLoTh, log.RadarCarCalHiTh,
			log.RadarCarDeltaTh, log.DownlinkEn7BitsRepeatedOccupancyPeriodMins,
		)
	}

	// Construct the SQL statement by joining the placeholders for each record
	query := fmt.Sprintf(
		"INSERT INTO %s (raw_id, device_id, firmware_version, network_type, happened_at, created_at, timestamp, device_mode, device_enable, radar_car_cal_lo_th, radar_car_cal_hi_th, radar_car_delta_th, downlink_en_7_bits_repeated_occupancy_period_mins) VALUES %s",
		s.TableName(), strings.Join(values, ", "),
	)

	// Execute the constructed query with the arguments
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk insert for setting logs: %w", err)
	}

	return nil
}
