package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/google/uuid"
)

// NbiotSettingLog represents a single NB-IoT setting log entry.
type NbiotSettingLog struct {
	ID         int       `db:"id" json:"id"`                   // Auto-incrementing primary key
	RawID      uuid.UUID `db:"raw_id" json:"raw_id"`           // ID linking to raw data source
	DeviceID   string    `db:"device_id" json:"device_id"`     // Device identifier, can be IMEI or UUID
	HappenedAt time.Time `db:"happened_at" json:"happened_at"` // Time when the setting event occurred
	CreatedAt  time.Time `db:"created_at" json:"created_at"`   // Time when the record was created
	Timestamp  int64     `db:"timestamp" json:"timestamp"`     // Event timestamp in UNIX format

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
	NBIoTIMSI                   int64  `db:"nb_iot_imsi" json:"nb_iot_imsi"`
}

// TableName returns the table name for the NbiotSettingLog model.
func (n *NbiotSettingLog) TableName() string {
	return "parking.nbiot_setting_logs"
}

// NewNbiotSettingLog constructs an NbiotSettingLog object from a provided map of data.
// It handles data type conversions and populates the fields accordingly.
func NewNbiotSettingLog(pktData map[string]any) (*NbiotSettingLog, error) {
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

	// Construct and return the NbiotSettingLog object with the parsed and converted data.
	nbiotSettingLog := NbiotSettingLog{
		RawID:             rawUUID,
		DeviceID:          pktData["device_id"].(string),
		HappenedAt:        happenedAt,
		CreatedAt:         time.Now().UTC(), // Default to the current time in UTC.
		Timestamp:         timestamp,
		DeviceMode:        int(pktData["device_mode"].(float64)),
		DeviceEnable:      int(pktData["device_enable"].(float64)),
		RadarCarCalLoTh:   int(pktData["radar_car_cal_lo_th"].(float64)),
		RadarCarCalHiTh:   int(pktData["radar_car_cal_hi_th"].(float64)),
		RadarCarUncalLoTh: int(pktData["radar_car_uncal_lo_th"].(float64)),
		RadarCarUncalHiTh: int(pktData["radar_car_uncal_hi_th"].(float64)),
		RadarCarDeltaTh:   int(pktData["radar_car_delta_th"].(float64)),
		MagCarLo:          int(pktData["mag_car_lo"].(float64)),
		MagCarHi:          int(pktData["mag_car_hi"].(float64)),

		// RadarTrailCalLoTh: int(pktData["radar_trail_cal_lo_th"].(float64)),
		// RadarTrailCalHiTh: int(pktData["radar_trail_cal_hi_th"].(float64)),
		// RadarTrailUncalLoTh: int(pktData["radar_trail_uncal_lo_th"].(float64)),
		// RadarTrailUncalHiTh: int(pktData["radar_trail_uncal_hi_th"].(float64)),

		DebugPeriod:                 int(pktData["debug_period"].(float64)),
		DebugMode:                   int(pktData["debug_mode"].(float64)),
		LogsMode:                    int(pktData["logs_mode"].(float64)),
		LogsAmount:                  int(pktData["logs_amount"].(float64)),
		MaximumRegistrationTime:     int(pktData["maximum_registration_time"].(float64)),
		MaximumRegistrationAttempts: int(pktData["maximum_registration_attempts"].(float64)),
		MaximumDeepSleepTime:        int(pktData["maximum_deep_sleep_time"].(float64)),
		TenXDeepSleepTime:           int64(pktData["ten_x_deep_sleep_time"].(float64)),
		TenXActionBefore:            int64(pktData["ten_x_action_before"].(float64)),
		TenXActionAfter:             int64(pktData["ten_x_action_after"].(float64)),
		NBIoTUDPIP:                  pktData["nb_iot_udp_ip"].(string),
		NBIoTUDPPort:                int(pktData["nb_iot_udp_port"].(float64)),
		NBIoTAPNLength:              int(pktData["nb_iot_apn_length"].(float64)),
		NBIoTAPN:                    pktData["nb_iot_apn"].(string),
		NBIoTIMSI:                   int64(pktData["nb_iot_imsi"].(float64)),
	}

	if pktData["radar_trail_cal_lo_th"] != nil {
		nbiotSettingLog.RadarTrailCalLoTh = int(pktData["radar_trail_cal_lo_th"].(float64))
	}
	if pktData["radar_trail_cal_hi_th"] != nil {
		nbiotSettingLog.RadarTrailCalHiTh = int(pktData["radar_trail_cal_hi_th"].(float64))
	}
	if pktData["radar_trail_uncal_lo_th"] != nil {
		nbiotSettingLog.RadarTrailUncalLoTh = int(pktData["radar_trail_uncal_lo_th"].(float64))
	}
	if pktData["radar_trail_uncal_hi_th"] != nil {
		nbiotSettingLog.RadarTrailUncalHiTh = int(pktData["radar_trail_uncal_hi_th"].(float64))
	}

	return &nbiotSettingLog, nil
}
