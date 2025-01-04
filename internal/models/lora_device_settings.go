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
