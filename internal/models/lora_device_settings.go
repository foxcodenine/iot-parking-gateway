package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// LoraDeviceSettings represents current or last settings of a LoRa device.
type LoraDeviceSettings struct {
	DeviceID        string    `db:"device_id" json:"device_id" primary_key:"true"` // Device identifier, serves as primary key
	FirmwareVersion float64   `db:"firmware_version" json:"firmware_version"`      // Firmware version of the device
	NetworkType     string    `db:"network_type" json:"network_type"`              // Network type (e.g., LoRa)
	CreatedAt       time.Time `db:"created_at" json:"created_at"`                  // Time when the record was created
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`                  // Time when the record was updated
	Timestamp       int64     `db:"timestamp" json:"timestamp"`                    // Event timestamp in UNIX format
	Flag            int       `db:"flag" json:"flag" default:"0"`                  // Additional flag column with a default value of 0

	DeviceMode                  int `db:"device_mode" json:"device_mode"`
	DeviceEnable                int `db:"device_enable" json:"device_enable"`
	RadarCarCalLoTh             int `db:"radar_car_cal_lo_th" json:"radar_car_cal_lo_th"`
	RadarCarCalHiTh             int `db:"radar_car_cal_hi_th" json:"radar_car_cal_hi_th"`
	RadarCarUncalLoTh           int `db:"radar_car_uncal_lo_th" json:"radar_car_uncal_lo_th"`
	RadarCarUncalHiTh           int `db:"radar_car_uncal_hi_th" json:"radar_car_uncal_hi_th"`
	RadarCarDeltaTh             int `db:"radar_car_delta_th" json:"radar_car_delta_th"`
	MagCarLo                    int `db:"mag_car_lo" json:"mag_car_lo"`
	MagCarHi                    int `db:"mag_car_hi" json:"mag_car_hi"`
	DebugPeriod                 int `db:"debug_period" json:"debug_period"`
	DebugMode                   int `db:"debug_mode" json:"debug_mode"`
	LogsMode                    int `db:"logs_mode" json:"logs_mode"`
	LogsAmount                  int `db:"logs_amount" json:"logs_amount"`
	MaximumRegistrationTime     int `db:"maximum_registration_time" json:"maximum_registration_time"`
	MaximumRegistrationAttempts int `db:"maximum_registration_attempts" json:"maximum_registration_attempts"`
	MaximumDeepSleepTime        int `db:"maximum_deep_sleep_time" json:"maximum_deep_sleep_time"`

	DeepSleepTime1  int `db:"deep_sleep_time_1" json:"deep_sleep_time_1"`
	ActionBefore1   int `db:"action_before_1" json:"action_before_1"`
	ActionAfter1    int `db:"action_after_1" json:"action_after_1"`
	DeepSleepTime2  int `db:"deep_sleep_time_2" json:"deep_sleep_time_2"`
	ActionBefore2   int `db:"action_before_2" json:"action_before_2"`
	ActionAfter2    int `db:"action_after_2" json:"action_after_2"`
	DeepSleepTime3  int `db:"deep_sleep_time_3" json:"deep_sleep_time_3"`
	ActionBefore3   int `db:"action_before_3" json:"action_before_3"`
	ActionAfter3    int `db:"action_after_3" json:"action_after_3"`
	DeepSleepTime4  int `db:"deep_sleep_time_4" json:"deep_sleep_time_4"`
	ActionBefore4   int `db:"action_before_4" json:"action_before_4"`
	ActionAfter4    int `db:"action_after_4" json:"action_after_4"`
	DeepSleepTime5  int `db:"deep_sleep_time_5" json:"deep_sleep_time_5"`
	ActionBefore5   int `db:"action_before_5" json:"action_before_5"`
	ActionAfter5    int `db:"action_after_5" json:"action_after_5"`
	DeepSleepTime6  int `db:"deep_sleep_time_6" json:"deep_sleep_time_6"`
	ActionBefore6   int `db:"action_before_6" json:"action_before_6"`
	ActionAfter6    int `db:"action_after_6" json:"action_after_6"`
	DeepSleepTime7  int `db:"deep_sleep_time_7" json:"deep_sleep_time_7"`
	ActionBefore7   int `db:"action_before_7" json:"action_before_7"`
	ActionAfter7    int `db:"action_after_7" json:"action_after_7"`
	DeepSleepTime8  int `db:"deep_sleep_time_8" json:"deep_sleep_time_8"`
	ActionBefore8   int `db:"action_before_8" json:"action_before_8"`
	ActionAfter8    int `db:"action_after_8" json:"action_after_8"`
	DeepSleepTime9  int `db:"deep_sleep_time_9" json:"deep_sleep_time_9"`
	ActionBefore9   int `db:"action_before_9" json:"action_before_9"`
	ActionAfter9    int `db:"action_after_9" json:"action_after_9"`
	DeepSleepTime10 int `db:"deep_sleep_time_10" json:"deep_sleep_time_10"`
	ActionBefore10  int `db:"action_before_10" json:"action_before_10"`
	ActionAfter10   int `db:"action_after_10" json:"action_after_10"`

	LoraDataRate int `db:"lora_data_rate" json:"lora_data_rate"`
	LoraRetries  int `db:"lora_retries" json:"lora_retries"`
}

// TableName returns the table name for the LoraDeviceSettings model.
func (l *LoraDeviceSettings) TableName() string {
	return "parking.lora_device_settings"
}

// Create inserts a new LoraDeviceSettings record into the database.
func (l *LoraDeviceSettings) Create(newSettings *LoraDeviceSettings) (*LoraDeviceSettings, error) {
	collection := dbSession.Collection(l.TableName())

	// Set the current time for CreatedAt
	newSettings.CreatedAt = time.Now().UTC()
	newSettings.UpdatedAt = time.Now().UTC()

	// Attempt to insert the new record
	_, err := collection.Insert(newSettings)
	if err != nil {
		// Check if the error is due to a duplicate key, which would imply a record with the same DeviceID already exists
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return nil, errors.New("a device settings with this ID already exists")
		}
		// Return any other errors with additional context
		return nil, fmt.Errorf("failed to create device settings: %w", err)
	}

	return newSettings, nil
}

// BulkUpdate updates multiple LoraDeviceSettings records based on their DeviceID.
func (l *LoraDeviceSettings) BulkUpdate(settings []LoraSettingLog) error {
	if len(settings) == 0 {
		return nil // No data to update
	}

	var args []interface{}
	valuesList := make([]string, len(settings))

	// Prepare the VALUES clause and arguments for the update query
	for i, setting := range settings {
		pos := i*52 + 1 // Adjust the position offset based on the number of fields per record
		valuesList[i] = fmt.Sprintf(
			"($%d,  $%d::numeric, $%d, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric, $%d::numeric)",
			pos, pos+1, pos+2, pos+3, pos+4, pos+5, pos+6, pos+7, pos+8, pos+9,
			pos+10, pos+11, pos+12, pos+13, pos+14, pos+15, pos+16, pos+17, pos+18, pos+19,
			pos+20, pos+21, pos+22, pos+23, pos+24, pos+25, pos+26, pos+27, pos+28, pos+29,
			pos+30, pos+31, pos+32, pos+33, pos+34, pos+35, pos+36, pos+37, pos+38, pos+39,
			pos+40, pos+41, pos+42, pos+43, pos+44, pos+45, pos+46, pos+47, pos+48, pos+49,
			pos+50, pos+51,
		)
		args = append(args,
			setting.DeviceID, setting.FirmwareVersion, setting.NetworkType, setting.Timestamp,
			setting.DeviceMode, setting.DeviceEnable, setting.RadarCarCalLoTh, setting.RadarCarCalHiTh,
			setting.RadarCarUncalLoTh, setting.RadarCarUncalHiTh, setting.RadarCarDeltaTh, setting.MagCarLo,
			setting.MagCarHi, setting.DebugPeriod, setting.DebugMode, setting.LogsMode, setting.LogsAmount,
			setting.MaximumRegistrationTime, setting.MaximumRegistrationAttempts, setting.MaximumDeepSleepTime,

			setting.DeepSleepTime1, setting.ActionBefore1, setting.ActionAfter1,
			setting.DeepSleepTime2, setting.ActionBefore2, setting.ActionAfter2,
			setting.DeepSleepTime3, setting.ActionBefore3, setting.ActionAfter3,
			setting.DeepSleepTime4, setting.ActionBefore4, setting.ActionAfter4,
			setting.DeepSleepTime5, setting.ActionBefore5, setting.ActionAfter5,
			setting.DeepSleepTime6, setting.ActionBefore6, setting.ActionAfter6,
			setting.DeepSleepTime7, setting.ActionBefore7, setting.ActionAfter7,
			setting.DeepSleepTime8, setting.ActionBefore8, setting.ActionAfter8,
			setting.DeepSleepTime9, setting.ActionBefore9, setting.ActionAfter9,
			setting.DeepSleepTime10, setting.ActionBefore10, setting.ActionAfter10,

			setting.LoraDataRate,
			setting.LoraRetries,
		)
	}

	// Construct the SQL query
	query := fmt.Sprintf(`
		UPDATE parking.lora_device_settings AS s
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
			debug_period = v.debug_period,
			debug_mode = v.debug_mode,
			logs_mode = v.logs_mode,
			logs_amount = v.logs_amount,
			maximum_registration_time = v.maximum_registration_time,
			maximum_registration_attempts = v.maximum_registration_attempts,
			maximum_deep_sleep_time = v.maximum_deep_sleep_time,
			deep_sleep_time_1 = v.deep_sleep_time_1,
			action_before_1 = v.action_before_1,
			action_after_1 = v.action_after_1,
			deep_sleep_time_2 = v.deep_sleep_time_2,
			action_before_2 = v.action_before_2,
			action_after_2 = v.action_after_2,
			deep_sleep_time_3 = v.deep_sleep_time_3,
			action_before_3 = v.action_before_3,
			action_after_3 = v.action_after_3,
			deep_sleep_time_4 = v.deep_sleep_time_4,
			action_before_4 = v.action_before_4,
			action_after_4 = v.action_after_4,
			deep_sleep_time_5 = v.deep_sleep_time_5,
			action_before_5 = v.action_before_5,
			action_after_5 = v.action_after_5,
			deep_sleep_time_6 = v.deep_sleep_time_6,
			action_before_6 = v.action_before_6,
			action_after_6 = v.action_after_6,
			deep_sleep_time_7 = v.deep_sleep_time_7,
			action_before_7 = v.action_before_7,
			action_after_7 = v.action_after_7,
			deep_sleep_time_8 = v.deep_sleep_time_8,
			action_before_8 = v.action_before_8,
			action_after_8 = v.action_after_8,
			deep_sleep_time_9 = v.deep_sleep_time_9,
			action_before_9 = v.action_before_9,
			action_after_9 = v.action_after_9,
			deep_sleep_time_10 = v.deep_sleep_time_10,
			action_before_10 = v.action_before_10,
			action_after_10 = v.action_after_10,			
			lora_data_rate = v.lora_data_rate,
			lora_retries = v.lora_retries

	
		FROM (VALUES %s) AS v(
			device_id, firmware_version, network_type, timestamp,
			device_mode, device_enable, radar_car_cal_lo_th, radar_car_cal_hi_th,
			radar_car_uncal_lo_th, radar_car_uncal_hi_th, radar_car_delta_th,
			mag_car_lo, mag_car_hi, debug_period, debug_mode, logs_mode, logs_amount,
			maximum_registration_time, maximum_registration_attempts, maximum_deep_sleep_time,
			deep_sleep_time_1, action_before_1, action_after_1, 
			deep_sleep_time_2, action_before_2, action_after_2, 
			deep_sleep_time_3, action_before_3, action_after_3, 
			deep_sleep_time_4, action_before_4, action_after_4, 
			deep_sleep_time_5, action_before_5, action_after_5, 
			deep_sleep_time_6, action_before_6, action_after_6, 
			deep_sleep_time_7, action_before_7, action_after_7, 
			deep_sleep_time_8, action_before_8, action_after_8, 
			deep_sleep_time_9, action_before_9, action_after_9, 
			deep_sleep_time_10, action_before_10, action_after_10, 
			lora_data_rate, lora_retries
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
