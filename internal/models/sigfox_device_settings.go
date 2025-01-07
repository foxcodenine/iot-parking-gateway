package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// SigfoxDeviceSettings represents the current or last settings of a Sigfox device.
type SigfoxDeviceSettings struct {
	DeviceID        string    `db:"device_id" json:"device_id" primary_key:"true"` // Device identifier, serves as primary key
	FirmwareVersion float64   `db:"firmware_version" json:"firmware_version"`      // Firmware version of the device
	NetworkType     string    `db:"network_type" json:"network_type"`              // Network type (e.g., Sigfox)
	CreatedAt       time.Time `db:"created_at" json:"created_at"`                  // Time when the record was created
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`                  // Time when the record was updated
	Timestamp       int64     `db:"timestamp" json:"timestamp"`                    // Event timestamp in UNIX format
	Flag            int       `db:"flag" json:"flag" default:"0"`                  // Additional flag column with a default value of 0

	DeviceMode                                 int `db:"device_mode" json:"device_mode"`                                                                             // Device mode
	DeviceEnable                               int `db:"device_enable" json:"device_enable"`                                                                         // Device enable status
	RadarCarCalLoTh                            int `db:"radar_car_cal_lo_th" json:"radar_car_cal_lo_th"`                                                             // Radar car calibration low threshold
	RadarCarCalHiTh                            int `db:"radar_car_cal_hi_th" json:"radar_car_cal_hi_th"`                                                             // Radar car calibration high threshold
	RadarCarDeltaTh                            int `db:"radar_car_delta_th" json:"radar_car_delta_th"`                                                               // Radar car delta threshold
	DownlinkEn7BitsRepeatedOccupancyPeriodMins int `db:"downlink_en_7_bits_repeated_occupancy_period_mins" json:"downlink_en_7_bits_repeated_occupancy_period_mins"` // Downlink period for 7 bits repeated occupancy
}

// TableName returns the table name for the SigfoxDeviceSettings model.
func (s *SigfoxDeviceSettings) TableName() string {
	return "parking.sigfox_device_settings"
}

// Create inserts a new SigfoxDeviceSettings record into the database.
func (s *SigfoxDeviceSettings) Create(newSettings *SigfoxDeviceSettings) (*SigfoxDeviceSettings, error) {
	// Get the database collection for SigfoxDeviceSettings
	collection := dbSession.Collection(s.TableName())

	// Set the current time for CreatedAt and UpdatedAt
	newSettings.CreatedAt = time.Now().UTC()
	newSettings.UpdatedAt = time.Now().UTC()

	// Attempt to insert the new record
	_, err := collection.Insert(newSettings)
	if err != nil {
		// Check if the error is due to a duplicate key (e.g., existing DeviceID)
		if strings.Contains(err.Error(), "SQLSTATE 23505") { // PostgreSQL unique constraint violation
			return nil, errors.New("a device settings with this ID already exists")
		}
		// Return any other errors with additional context
		return nil, fmt.Errorf("failed to create Sigfox device settings: %w", err)
	}

	return newSettings, nil
}

// BulkUpdate updates multiple SigfoxDeviceSettings records based on their DeviceID.
func (s *SigfoxDeviceSettings) BulkUpdate(settings []SigfoxSettingLog) error {
	if len(settings) == 0 {
		return nil // No data to update
	}

	var args []interface{}
	valuesList := make([]string, len(settings))

	// Prepare the VALUES clause and arguments for the update query
	for i, setting := range settings {
		pos := i*10 + 1 // Adjust the position offset based on the number of fields per record
		valuesList[i] = fmt.Sprintf(
			"($%d, $%d::numeric, $%d, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric)",
			pos, pos+1, pos+2, pos+3, pos+4, pos+5, pos+6, pos+7, pos+8, pos+9,
		)
		args = append(args,
			setting.DeviceID, setting.FirmwareVersion, setting.NetworkType, setting.Timestamp,
			setting.DeviceMode, setting.DeviceEnable, setting.RadarCarCalLoTh, setting.RadarCarCalHiTh,
			setting.RadarCarDeltaTh, setting.DownlinkEn7BitsRepeatedOccupancyPeriodMins,
		)
	}

	// Construct the SQL query
	query := fmt.Sprintf(`
		UPDATE parking.sigfox_device_settings AS s
		SET
			firmware_version = v.firmware_version,
			network_type = v.network_type,
			timestamp = v.timestamp,
			device_mode = v.device_mode,
			device_enable = v.device_enable,
			radar_car_cal_lo_th = v.radar_car_cal_lo_th,
			radar_car_cal_hi_th = v.radar_car_cal_hi_th,
			radar_car_delta_th = v.radar_car_delta_th,
			downlink_en_7_bits_repeated_occupancy_period_mins = v.downlink_en_7_bits_repeated_occupancy_period_mins
		FROM (VALUES %s) AS v(
			device_id, firmware_version, network_type, timestamp,
			device_mode, device_enable, radar_car_cal_lo_th, radar_car_cal_hi_th,
			radar_car_delta_th, downlink_en_7_bits_repeated_occupancy_period_mins
		)
		WHERE s.device_id = v.device_id
	`, strings.Join(valuesList, ", "))

	// Execute the query
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk update for Sigfox device settings: %w", err)
	}

	return nil
}
