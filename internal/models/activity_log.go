package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	up "github.com/upper/db/v4"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/google/uuid"
)

// ActivityLog struct to store activity logs
type ActivityLog struct {
	ID              int          `db:"id" json:"id"`                             // Auto-incrementing primary key
	RawID           uuid.UUID    `db:"raw_id" json:"raw_id"`                     // ID linking to raw data source
	DeviceID        string       `db:"device_id" json:"device_id"`               // Device identifier, can be IMEI or UUID
	FirmwareVersion float64      `db:"firmware_version" json:"firmware_version"` // Firmware version of the device
	NetworkType     string       `db:"network_type" json:"network_type"`         // Network type (e.g., NB-IoT, LoRa, Sigfox)
	HappenedAt      time.Time    `db:"happened_at" json:"happened_at"`           // Time when the activity event occurred
	CreatedAt       time.Time    `db:"created_at" json:"created_at"`             // Time when the record was created
	Timestamp       int64        `db:"timestamp" json:"timestamp"`               // Epoch time for additional timing information
	BeaconsAmount   int          `db:"beacons_amount" json:"beacons_amount"`     // Number of beacons involved
	MagnetAbsTotal  int          `db:"magnet_abs_total" json:"magnet_abs_total"` // Total magnetic reading value
	PeakDistanceCm  int          `db:"peak_distance_cm" json:"peak_distance_cm"` // Peak distance in centimeters
	RadarCumulative int          `db:"radar_cumulative" json:"radar_cumulative"` // Cumulative radar reading
	IsOccupied      bool         `db:"is_occupied" json:"is_occupied"`           // Whether a vehicle is detected
	Beacons         *BeaconSlice `db:"beacons" json:"beacons"`                   // JSONB column to store an array of beacon data
}

// TableName returns the table name for the ActivityLog model.
func (a *ActivityLog) TableName() string {
	return "parking.activity_logs"
}

// NewActivityLog constructs an ActivityLog object from a provided map of data.
// It validates and converts the data types as necessary, handling potential errors at each step.
func NewActivityLog(pktData map[string]any) (*ActivityLog, error) {
	// Extract the 'raw_id' from pktData and validate it's a string.
	rawUUIDStr, ok := pktData["raw_id"].(string)
	if !ok {
		return nil, errors.New("invalid uuid format: expected string")
	}

	// Attempt to parse the 'raw_id' into a UUID format.
	rawUUID, err := uuid.Parse(rawUUIDStr)
	if err != nil {
		// Wrap the error to provide context if UUID parsing fails.
		return nil, helpers.WrapError(fmt.Errorf("failed to parse uuid %s: %v", rawUUIDStr, err))
	}

	// Extract the 'timestamp' field and validate it's a float64.
	timestampFloat, ok := pktData["timestamp"].(float64)
	if !ok {
		return nil, errors.New("invalid timestamp format: expected float64")
	}

	// Convert the timestamp from float64 to int64 and then to time.Time in UTC.
	timestampInt := int64(timestampFloat)
	happenedAt := time.Unix(timestampInt, 0).UTC()

	// Initialize a variable to hold beacon data.
	var beacons BeaconSlice

	// Extract the 'beacons' field and validate it's an array of interfaces (generic JSON array).
	beaconData, ok := pktData["beacons"].([]interface{})
	if !ok {
		return nil, errors.New("invalid type for beacons: expected []interface{}")
	}

	// Marshal the generic interface array back to JSON bytes.
	beaconBytes, err := json.Marshal(beaconData)
	if err != nil {
		return nil, errors.New("error marshalling beacons data")
	}

	// Unmarshal the JSON bytes back into the BeaconSlice type.
	err = json.Unmarshal(beaconBytes, &beacons)
	if err != nil {
		return nil, errors.New("error unmarshalling beacons data")
	}

	// Extract and validate the 'network_type' field.
	networkType, ok := pktData["network_type"].(string)
	if !ok {
		return nil, errors.New("invalid type for network_type: expected string")
	}

	// Create the ActivityLog object using the extracted and converted fields.
	activityLog := &ActivityLog{
		RawID:           rawUUID,
		DeviceID:        pktData["device_id"].(string),
		FirmwareVersion: pktData["firmware_version"].(float64),
		NetworkType:     networkType,
		HappenedAt:      happenedAt,
		Timestamp:       timestampInt,
		BeaconsAmount:   int(pktData["beacons_amount"].(float64)),
		// MagnetAbsTotal:  int(pktData["magnet_abs_total"].(float64)),
		// PeakDistanceCm:  int(pktData["peak_distance_cm"].(float64)),
		RadarCumulative: int(pktData["radar_cumulative"].(float64)),
		IsOccupied:      pktData["is_occupied"].(float64) != 0, // Convert float64 to boolean (non-zero is true).
		Beacons:         &beacons,                              // Attach the processed beacon slice.
	}

	peakDistanceCm, ok := pktData["peak_distance_cm"]
	if ok {
		activityLog.PeakDistanceCm = int(peakDistanceCm.(float64))
	}

	// Check for the optional "magnet_abs_total" field.
	if val, ok := pktData["magnet_abs_total"].(float64); ok {
		activityLog.MagnetAbsTotal = int(val)
	}

	// Return the newly created ActivityLog object and nil as the error.
	return activityLog, nil
}

// BulkInsert inserts multiple ActivityLog records in a single operation.
func (a *ActivityLog) BulkInsert(activityLogs []ActivityLog) error {
	// Exit early if there are no records to insert
	if len(activityLogs) == 0 {
		return nil
	}

	// Prepare slices for SQL values and arguments.
	values := make([]string, 0, len(activityLogs))       // Holds the placeholder for each row
	args := make([]interface{}, 0, len(activityLogs)*12) // Updated argument count to include firmware_version

	for i, log := range activityLogs {

		// Create a placeholder for each record with indexed arguments, e.g., ($1, $2, ..., $12)
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*12+1, i*12+2, i*12+3, i*12+4, i*12+5, i*12+6, i*12+7, i*12+8, i*12+9, i*12+10, i*12+11, i*12+12))

		// Append the actual values for each placeholder in the same order as the columns
		args = append(args, log.RawID, log.DeviceID, log.FirmwareVersion, log.NetworkType, log.HappenedAt, log.Timestamp,
			log.BeaconsAmount, log.MagnetAbsTotal, log.PeakDistanceCm, log.RadarCumulative, log.IsOccupied, log.Beacons)
	}

	// Construct the SQL statement by joining the placeholders for each record
	query := fmt.Sprintf("INSERT INTO %s (raw_id, device_id, firmware_version, network_type, happened_at, timestamp, beacons_amount, magnet_abs_total, peak_distance_cm, radar_cumulative, is_occupied, beacons) VALUES %s",
		a.TableName(), strings.Join(values, ", "))

	// Execute the constructed query with the arguments
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return helpers.WrapError(fmt.Errorf("failed to execute bulk insert for activity logs: %w", err))
	}

	return nil
}

func (a *ActivityLog) GetLastEntryFromPreviousDay() (*ActivityLog, error) {
	// Use UTC for consistency; adjust if you need a different timezone.
	now := time.Now().UTC()

	// Calculate the start of today (00:00:00) and the start of yesterday.
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startOfYesterday := startOfToday.AddDate(0, 0, -1)

	collection := dbSession.Collection(a.TableName())

	fmt.Println(startOfYesterday, startOfToday)

	var lastLog ActivityLog
	err := collection.
		Find("happened_at >= ? AND happened_at < ?", startOfYesterday, startOfToday).
		OrderBy("happened_at DESC").
		One(&lastLog)
	if err != nil {
		if err == up.ErrNoMoreRows {
			return nil, nil
		}
		return nil, helpers.WrapError(fmt.Errorf("failed to retrieve last activity log from previous day: %w", err))
	}

	fmt.Println(lastLog)

	return &lastLog, nil
}

func (a *ActivityLog) GetActivityLogs(deviceID string, fromDate, toDate int64) ([]*ActivityLog, error) {
	// Validate inputs
	if deviceID == "" {
		return nil, helpers.WrapError(fmt.Errorf("device_id cannot be empty"))
	}
	if fromDate <= 0 || toDate <= 0 {
		return nil, helpers.WrapError(fmt.Errorf("invalid timestamps: both fromDate and toDate must be greater than zero"))
	}
	if fromDate > toDate {
		return nil, helpers.WrapError(fmt.Errorf("fromDate cannot be greater than toDate"))
	}

	// Convert timestamps to time.Time
	fromTime := time.Unix(fromDate, 0).UTC()
	toTime := time.Unix(toDate, 0).UTC()

	// Get the collection
	collection := dbSession.Collection(a.TableName())

	// Query the database for logs in the time range
	var logs []*ActivityLog

	err := collection.
		Find("device_id = ? AND happened_at BETWEEN ? AND ?", deviceID, fromTime, toTime).
		OrderBy("happened_at ASC").
		All(&logs)

	// Handle query errors
	if err != nil {
		if err == up.ErrNoMoreRows {
			return nil, nil // No logs found
		}
		return nil, helpers.WrapError(fmt.Errorf("failed to fetch activity logs: %w", err))
	}

	return logs, nil
}
