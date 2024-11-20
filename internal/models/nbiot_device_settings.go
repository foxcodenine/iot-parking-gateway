package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// NbiotDeviceSettings represents current or last settings of a device.
type NbiotDeviceSettings struct {
	DeviceID        string    `db:"device_id" json:"device_id" primary_key:"true"` // Device identifier, serves as primary key
	FirmwareVersion float64   `db:"firmware_version" json:"firmware_version"`      // Firmware version of the device
	NetworkType     string    `db:"network_type" json:"network_type"`              // Network type (e.g., NB-IoT)
	CreatedAt       time.Time `db:"created_at" json:"created_at"`                  // Time when the record was created
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`                  // Time when the record was updated
	Timestamp       int64     `db:"timestamp" json:"timestamp"`                    // Event timestamp in UNIX format
	Flag            int       `db:"flag" json:"flag" default:"0"`                  // Additional flag column with a default value of 0

	DeviceMode                  int    `db:"device_mode" json:"device_mode"`
	DeviceEnable                int    `db:"device_enable" json:"device_enable"`
	RadarCarCalLoTh             int    `db:"radar_car_cal_lo_th" json:"radar_car_cal_lo_th"`
	RadarCarCalHiTh             int    `db:"radar_car_cal_hi_th" json:"radar_car_cal_hi_th"`
	RadarCarUncalLoTh           int    `db:"radar_car_uncal_lo_th" json:"radar_car_uncal_lo_th"`
	RadarCarUncalHiTh           int    `db:"radar_car_uncal_hi_th" json:"radar_car_uncal_hi_th"`
	RadarCarDeltaTh             int    `db:"radar_car_delta_th" json:"radar_car_delta_th"`
	MagCarLo                    int    `db:"mag_car_lo" json:"mag_car_lo"`
	MagCarHi                    int    `db:"mag_car_hi" json:"mag_car_hi"`
	RadarTrailCalLoTh           int    `db:"radar_trail_cal_lo_th" json:"radar_trail_cal_lo_th"`
	RadarTrailCalHiTh           int    `db:"radar_trail_cal_hi_th" json:"radar_trail_cal_hi_th"`
	RadarTrailUncalLoTh         int    `db:"radar_trail_uncal_lo_th" json:"radar_trail_uncal_lo_th"`
	RadarTrailUncalHiTh         int    `db:"radar_trail_uncal_hi_th" json:"radar_trail_uncal_hi_th"`
	DebugPeriod                 int    `db:"debug_period" json:"debug_period"`
	DebugMode                   int    `db:"debug_mode" json:"debug_mode"`
	LogsMode                    int    `db:"logs_mode" json:"logs_mode"`
	LogsAmount                  int    `db:"logs_amount" json:"logs_amount"`
	MaximumRegistrationTime     int    `db:"maximum_registration_time" json:"maximum_registration_time"`
	MaximumRegistrationAttempts int    `db:"maximum_registration_attempts" json:"maximum_registration_attempts"`
	MaximumDeepSleepTime        int    `db:"maximum_deep_sleep_time" json:"maximum_deep_sleep_time"`
	TenXDeepSleepTime           int64  `db:"ten_x_deep_sleep_time" json:"ten_x_deep_sleep_time"`
	TenXActionBefore            int64  `db:"ten_x_action_before" json:"ten_x_action_before"`
	TenXActionAfter             int64  `db:"ten_x_action_after" json:"ten_x_action_after"`
	NBIoTUDPIP                  string `db:"nb_iot_udp_ip" json:"nb_iot_udp_ip"`
	NBIoTUDPPort                int    `db:"nb_iot_udp_port" json:"nb_iot_udp_port"`
	NBIoTAPNLength              int    `db:"nb_iot_apn_length" json:"nb_iot_apn_length"`
	NBIoTAPN                    string `db:"nb_iot_apn" json:"nb_iot_apn"`
	NBIoTIMSI                   string `db:"nb_iot_imsi" json:"nb_iot_imsi"`
}

// TableName returns the table name for the NbiotDeviceSettings model.
func (n *NbiotDeviceSettings) TableName() string {
	return "parking.nbiot_device_settings"
}

// Create inserts a new NbiotDeviceSettings record into the database.
func (n *NbiotDeviceSettings) Create(newSettings *NbiotDeviceSettings) (*NbiotDeviceSettings, error) {
	collection := dbSession.Collection(n.TableName())

	// Set the current time for CreatedAt
	newSettings.CreatedAt = time.Now().UTC()
	newSettings.UpdatedAt = time.Now().UTC()

	// Attempt to insert the new record
	_, err := collection.Insert(newSettings)
	if err != nil {
		// Check if the error is due to a duplicate key, which would imply a record with the same DeviceID already exists
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return nil, errors.New("a device with this ID already exists")
		}
		// Return any other errors with additional context
		return nil, fmt.Errorf("failed to create device settings: %w", err)
	}

	return newSettings, nil
}

// BulkUpdate updates multiple NbiotDeviceSettings records based on their DeviceID.
func (n *NbiotDeviceSettings) BulkUpdate(settings []NbiotSettingLog) error {
	if len(settings) == 0 {
		return nil // No data to update
	}

	var args []interface{}
	valuesList := make([]string, len(settings))

	// Prepare the VALUES clause and arguments for the update query
	for i, setting := range settings {
		pos := i*32 + 1 // Adjust the position offset based on the number of fields per record
		valuesList[i] = fmt.Sprintf(
			"($%d,  $%d::numeric, $%d, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d, $%d::numeric, $%d::numeric, $%d, $%d)",
			pos, pos+1, pos+2, pos+3, pos+4, pos+5, pos+6, pos+7, pos+8, pos+9,
			pos+10, pos+11, pos+12, pos+13, pos+14, pos+15, pos+16, pos+17, pos+18, pos+19,
			pos+20, pos+21, pos+22, pos+23, pos+24, pos+25, pos+26, pos+27, pos+28, pos+29, pos+30, pos+31,
		)
		args = append(args,
			setting.DeviceID, setting.FirmwareVersion, setting.NetworkType, setting.Timestamp,
			setting.DeviceMode, setting.DeviceEnable, setting.RadarCarCalLoTh, setting.RadarCarCalHiTh,
			setting.RadarCarUncalLoTh, setting.RadarCarUncalHiTh, setting.RadarCarDeltaTh, setting.MagCarLo,
			setting.MagCarHi, setting.DebugPeriod, setting.RadarTrailCalLoTh, setting.RadarTrailCalHiTh, setting.RadarTrailUncalLoTh, setting.RadarTrailUncalHiTh, setting.DebugMode, setting.LogsMode, setting.LogsAmount,
			setting.MaximumRegistrationTime, setting.MaximumRegistrationAttempts, setting.MaximumDeepSleepTime,
			setting.TenXDeepSleepTime, setting.TenXActionBefore, setting.TenXActionAfter, setting.NBIoTUDPIP,
			setting.NBIoTUDPPort, setting.NBIoTAPNLength, setting.NBIoTAPN, setting.NBIoTIMSI,
		)
	}

	// Construct the SQL query
	query := fmt.Sprintf(`
		UPDATE parking.nbiot_device_settings AS s
		SET
			firmware_version = v.firmware_version,
			network_type = v.network_type,
			timestamp = v.timestamp,		
			device_mode = v.device_mode,
			device_enable = v.device_enable,
			radar_car_cal_lo_th = v.radar_car_cal_lo_th,
			radar_car_cal_hi_th = v.radar_car_cal_hi_th,
			radar_car_uncal_lo_th = v.radar_car_uncal_lo_th,
			radar_car_uncal_hi_th = v.radar_car_uncal_hi_th,
			radar_car_delta_th = v.radar_car_delta_th,
			mag_car_lo = v.mag_car_lo,
			mag_car_hi = v.mag_car_hi,
			radar_trail_cal_lo_th = v.radar_trail_cal_lo_th,
			radar_trail_cal_hi_th = v.radar_trail_cal_hi_th,
			radar_trail_uncal_lo_th = v.radar_trail_uncal_lo_th,
			radar_trail_uncal_hi_th = v.radar_trail_uncal_hi_th,
			debug_period = v.debug_period,
			debug_mode = v.debug_mode,
			logs_mode = v.logs_mode,
			logs_amount = v.logs_amount,
			maximum_registration_time = v.maximum_registration_time,
			maximum_registration_attempts = v.maximum_registration_attempts,
			maximum_deep_sleep_time = v.maximum_deep_sleep_time,
			ten_x_deep_sleep_time = v.ten_x_deep_sleep_time,
			ten_x_action_before = v.ten_x_action_before,
			ten_x_action_after = v.ten_x_action_after,
			nb_iot_udp_ip = v.nb_iot_udp_ip,
			nb_iot_udp_port = v.nb_iot_udp_port,
			nb_iot_apn_length = v.nb_iot_apn_length,
			nb_iot_apn = v.nb_iot_apn,
			nb_iot_imsi = v.nb_iot_imsi
	
		FROM (VALUES %s) AS v(
			device_id, firmware_version, network_type, timestamp,
			device_mode, device_enable, radar_car_cal_lo_th, radar_car_cal_hi_th,
			radar_car_uncal_lo_th, radar_car_uncal_hi_th, radar_car_delta_th,
			mag_car_lo, mag_car_hi, radar_trail_cal_lo_th, radar_trail_cal_hi_th, radar_trail_uncal_lo_th, radar_trail_uncal_hi_th, debug_period, debug_mode, logs_mode, logs_amount,
			maximum_registration_time, maximum_registration_attempts, maximum_deep_sleep_time,
			ten_x_deep_sleep_time, ten_x_action_before, ten_x_action_after, nb_iot_udp_ip,
			nb_iot_udp_port, nb_iot_apn_length, nb_iot_apn, nb_iot_imsi
		)
		WHERE s.device_id = v.device_id
	`, strings.Join(valuesList, ", "))

	// Execute the query
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk update for device settings: %w", err)
	}

	return nil
}
