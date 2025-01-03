package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// LoraSettingLog represents a single LoRa setting log entry.
type LoraSettingLog struct {
	ID              int       `db:"id" json:"id"`
	RawID           uuid.UUID `db:"raw_id" json:"raw_id"`
	DeviceID        string    `db:"device_id" json:"device_id"`
	FirmwareVersion float64   `db:"firmware_version" json:"firmware_version"`
	NetworkType     string    `db:"network_type" json:"network_type"`
	HappenedAt      time.Time `db:"happened_at" json:"happened_at"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	Timestamp       int64     `db:"timestamp" json:"timestamp"`

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

// TableName returns the table name for the LoraSettingLog model.
func (l *LoraSettingLog) TableName() string {
	return "parking.lora_setting_logs"
}

// NewLoraSettingLog constructs a LoraSettingLog object from a provided map of data.
// It handles data type conversions and populates the fields accordingly.
func NewLoraSettingLog(pktData map[string]any) (*LoraSettingLog, error) {
	// Parse the raw UUID field from a string to uuid.UUID.
	rawUUIDStr, ok := pktData["raw_id"].(string)
	if !ok {
		return nil, errors.New("invalid uuid format: expected string")
	}

	rawUUID, err := uuid.Parse(rawUUIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid %s: %v", rawUUIDStr, err)
	}

	// Parse the timestamp, expected as an int64.
	timestampFloat, ok := pktData["timestamp"].(float64)
	if !ok {
		return nil, errors.New("invalid timestamp format: expected float64")
	}
	// Convert float64 timestamp to int64 and then to time.Time.
	timestamp := int64(timestampFloat)

	// Convert int64 timestamp to time.Time.
	happenedAt := time.Unix(timestamp, 0).UTC()

	// Construct and return the LoraSettingLog object with the parsed and converted data.
	log := &LoraSettingLog{
		RawID:           rawUUID,
		DeviceID:        pktData["device_id"].(string),
		FirmwareVersion: pktData["firmware_version"].(float64),
		NetworkType:     pktData["network_type"].(string),
		HappenedAt:      happenedAt,
		CreatedAt:       time.Now().UTC(), // Default to the current time in UTC.
		Timestamp:       timestamp,

		DeviceMode:                  int(pktData["device_mode"].(float64)),
		DeviceEnable:                int(pktData["device_enable"].(float64)),
		RadarCarCalLoTh:             int(pktData["radar_car_cal_lo_th"].(float64)),
		RadarCarCalHiTh:             int(pktData["radar_car_cal_hi_th"].(float64)),
		RadarCarUncalLoTh:           int(pktData["radar_car_uncal_lo_th"].(float64)),
		RadarCarUncalHiTh:           int(pktData["radar_car_uncal_hi_th"].(float64)),
		RadarCarDeltaTh:             int(pktData["radar_car_delta_th"].(float64)),
		MagCarLo:                    int(pktData["mag_car_lo"].(float64)),
		MagCarHi:                    int(pktData["mag_car_hi"].(float64)),
		DebugPeriod:                 int(pktData["debug_period"].(float64)),
		DebugMode:                   int(pktData["debug_mode"].(float64)),
		LogsMode:                    int(pktData["logs_mode"].(float64)),
		LogsAmount:                  int(pktData["logs_amount"].(float64)),
		MaximumRegistrationTime:     int(pktData["maximum_registration_time"].(float64)),
		MaximumRegistrationAttempts: int(pktData["maximum_registration_attempts"].(float64)),
		MaximumDeepSleepTime:        int(pktData["maximum_deep_sleep_time"].(float64)),

		DeepSleepTime1:  int(pktData["deep_sleep_time_1"].(float64)),
		ActionBefore1:   int(pktData["action_before_1"].(float64)),
		ActionAfter1:    int(pktData["action_after_1"].(float64)),
		DeepSleepTime2:  int(pktData["deep_sleep_time_2"].(float64)),
		ActionBefore2:   int(pktData["action_before_2"].(float64)),
		ActionAfter2:    int(pktData["action_after_2"].(float64)),
		DeepSleepTime3:  int(pktData["deep_sleep_time_3"].(float64)),
		ActionBefore3:   int(pktData["action_before_3"].(float64)),
		ActionAfter3:    int(pktData["action_after_3"].(float64)),
		DeepSleepTime4:  int(pktData["deep_sleep_time_4"].(float64)),
		ActionBefore4:   int(pktData["action_before_4"].(float64)),
		ActionAfter4:    int(pktData["action_after_4"].(float64)),
		DeepSleepTime5:  int(pktData["deep_sleep_time_5"].(float64)),
		ActionBefore5:   int(pktData["action_before_5"].(float64)),
		ActionAfter5:    int(pktData["action_after_5"].(float64)),
		DeepSleepTime6:  int(pktData["deep_sleep_time_6"].(float64)),
		ActionBefore6:   int(pktData["action_before_6"].(float64)),
		ActionAfter6:    int(pktData["action_after_6"].(float64)),
		DeepSleepTime7:  int(pktData["deep_sleep_time_7"].(float64)),
		ActionBefore7:   int(pktData["action_before_7"].(float64)),
		ActionAfter7:    int(pktData["action_after_7"].(float64)),
		DeepSleepTime8:  int(pktData["deep_sleep_time_8"].(float64)),
		ActionBefore8:   int(pktData["action_before_8"].(float64)),
		ActionAfter8:    int(pktData["action_after_8"].(float64)),
		DeepSleepTime9:  int(pktData["deep_sleep_time_9"].(float64)),
		ActionBefore9:   int(pktData["action_before_9"].(float64)),
		ActionAfter9:    int(pktData["action_after_9"].(float64)),
		DeepSleepTime10: int(pktData["deep_sleep_time_10"].(float64)),
		ActionBefore10:  int(pktData["action_before_10"].(float64)),
		ActionAfter10:   int(pktData["action_after_10"].(float64)),

		LoraDataRate: int(pktData["lora_data_rate"].(float64)),
		LoraRetries:  int(pktData["lora_retries"].(float64)),
	}

	return log, nil
}

// BulkInsert inserts multiple NbiotSettingLog records in a single operation.
func (l *LoraSettingLog) BulkInsert(settingLogs []NbiotSettingLog) error {
	// Exit early if there are no records to insert
	if len(settingLogs) == 0 {
		return nil
	}

	// Determine the number of fields to be inserted for each log
	numFields := 62 // Adjust this based on the actual number of columns you have in your table

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
			log.RadarCarUncalLoTh, log.RadarCarUncalHiTh, log.RadarCarDeltaTh, log.MagCarLo,
			log.MagCarHi, log.RadarTrailCalLoTh, log.RadarTrailCalHiTh, log.RadarTrailUncalLoTh,
			log.RadarTrailUncalHiTh, log.DebugPeriod, log.DebugMode, log.LogsMode, log.LogsAmount,
			log.MaximumRegistrationTime, log.MaximumRegistrationAttempts, log.MaximumDeepSleepTime,

			log.DeepSleepTime1, log.ActionBefore1, log.ActionAfter1,
			log.DeepSleepTime2, log.ActionBefore2, log.ActionAfter2,
			log.DeepSleepTime3, log.ActionBefore3, log.ActionAfter3,
			log.DeepSleepTime4, log.ActionBefore4, log.ActionAfter4,
			log.DeepSleepTime5, log.ActionBefore5, log.ActionAfter5,
			log.DeepSleepTime6, log.ActionBefore6, log.ActionAfter6,
			log.DeepSleepTime7, log.ActionBefore7, log.ActionAfter7,
			log.DeepSleepTime8, log.ActionBefore8, log.ActionAfter8,
			log.DeepSleepTime9, log.ActionBefore9, log.ActionAfter9,
			log.DeepSleepTime10, log.ActionBefore10, log.ActionAfter10,

			log.NBIoTUDPIP, log.NBIoTUDPPort, log.NBIoTAPNLength, log.NBIoTAPN, log.NBIoTIMSI,
		)
	}

	// Construct the SQL statement by joining the placeholders for each record
	query := fmt.Sprintf("INSERT INTO %s (raw_id, device_id, firmware_version, network_type, happened_at, created_at, timestamp, device_mode, device_enable, radar_car_cal_lo_th, radar_car_cal_hi_th, radar_car_uncal_lo_th, radar_car_uncal_hi_th, radar_car_delta_th, mag_car_lo, mag_car_hi, radar_trail_cal_lo_th, radar_trail_cal_hi_th, radar_trail_uncal_lo_th, radar_trail_uncal_hi_th, debug_period, debug_mode, logs_mode, logs_amount, maximum_registration_time, maximum_registration_attempts, maximum_deep_sleep_time, deep_sleep_time_1, action_before_1, action_after_1,	deep_sleep_time_2, action_before_2, action_after_2,	deep_sleep_time_3, action_before_3, action_after_3,	deep_sleep_time_4, action_before_4, action_after_4,	deep_sleep_time_5, action_before_5, action_after_5,	deep_sleep_time_6, action_before_6, action_after_6,	deep_sleep_time_7, action_before_7, action_after_7,	deep_sleep_time_8, action_before_8, action_after_8,	deep_sleep_time_9, action_before_9, action_after_9,	deep_sleep_time_10, action_before_10, action_after_10, nb_iot_udp_ip, nb_iot_udp_port, nb_iot_apn_length, nb_iot_apn, nb_iot_imsi) VALUES %s", l.TableName(), strings.Join(values, ", "))

	// Execute the constructed query with the arguments
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk insert for setting logs: %w", err)
	}

	return nil
}
