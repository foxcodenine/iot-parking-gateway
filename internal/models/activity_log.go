package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/google/uuid"
)

// ActivityLog
type ActivityLog struct {
	ID              int       `db:"id" json:"id"`                             // Auto-incrementing primary key
	RawID           uuid.UUID `db:"raw_id" json:"raw_id"`                     // ID linking to raw data source
	DeviceID        string    `db:"device_id" json:"device_id"`               // Device identifier, can be IMEI or UUID
	FirmwareVersion float64   `db:"firmware_version" json:"firmware_version"` // Firmware version of the device
	NetworkType     string    `db:"network_type" json:"network_type"`         // Network type (e.g., NB-IoT, LoRa, Sigfox)
	HappenedAt      time.Time `db:"happened_at" json:"happened_at"`           // Time when the activity event occurred
	CreatedAt       time.Time `db:"created_at" json:"created_at"`             // Time when the record was created
	Timestamp       int64     `db:"timestamp" json:"timestamp"`               // Epoch time for additional timing information
	BeaconsAmount   int       `db:"beacons_amount" json:"beacons_amount"`     // Number of beacons involved
	MagnetAbsTotal  int       `db:"magnet_abs_total" json:"magnet_abs_total"` // Total magnetic reading value
	PeakDistanceCm  int       `db:"peak_distance_cm" json:"peak_distance_cm"` // Peak distance in centimeters
	RadarCumulative int       `db:"radar_cumulative" json:"radar_cumulative"` // Cumulative radar reading
	IsOccupied      bool      `db:"is_occupied" json:"is_occupied"`           // Whether a vehicle is detected
	Beacons         []Beacon  `db:"beacons" json:"beacons"`                   // JSONB column to store an array of beacon data
}

// TableName returns the table name for the ActivityLog model.
func (a *ActivityLog) TableName() string {
	return "parking.activity_logs"
}

// NewActivityLog constructs an ActivityLog object from a provided map of data.
// It handles data type conversions and populates the fields accordingly.
func NewActivityLog(pktData map[string]any) (*ActivityLog, error) {

	// Convert the raw UUID field from a string to uuid.UUID.
	rawUUIDStr, ok := pktData["raw_id"].(string)
	if !ok {
		return nil, errors.New("invalid uuid format: expected string")
	}

	rawUUID, err := uuid.Parse(rawUUIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid %s: %v", rawUUIDStr, err)
	}

	// Parse the timestamp, expected as a float64 (representing seconds since epoch),
	// and convert it to time.Time.
	timestampFloat, ok := pktData["timestamp"].(float64)
	if !ok {
		helpers.LogError(nil, "invalid timestamp format: expected float64")
		return nil, errors.New("invalid timestamp format: expected float64")
	}

	// Convert float64 timestamp to int64 and then to time.Time.
	timestampInt := int64(timestampFloat)
	happenedAt := time.Unix(timestampInt, 0).UTC()

	// Initialize a slice to hold beacon entries.
	var beacons []Beacon

	// Retrieve and process the "beacons" data, expected as an array of maps.
	beaconData, ok := pktData["beacons"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type for beacons: expected []interface{}")
	}

	networkType, ok := pktData["network_type"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid type for network_type: expected string")
	}

	// Iterate over each beacon entry in the array.
	for _, beaconItem := range beaconData {
		// Assert that each beacon entry is a map with string keys and interface{} values.
		beaconMap, ok := beaconItem.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type for beacon item: expected map[string]interface{}")
		}

		// Safely retrieve and convert each beacon field from the map.
		beaconNumber, ok := beaconMap["beacon_number"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid type for beacon_number")
		}

		major, ok := beaconMap["major"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid type for major")
		}

		minor, ok := beaconMap["minor"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid type for minor")
		}

		// Append the beacon to the Beacons slice after converting values to int.
		beacon := Beacon{
			BeaconNumber: int(beaconNumber),
			Major:        int(major),
			Minor:        int(minor),
		}

		if networkType == "NB-Iot" || networkType == "LoRa" {
			rssi, ok := beaconMap["rssi"].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid type for rssi")
			}

			beacon.RSSI = int(rssi)
			beacons = append(beacons, beacon)
		}
	}

	// Construct and return the ActivityLog object with the parsed and converted data.
	activityLog := &ActivityLog{
		RawID:           rawUUID,
		DeviceID:        pktData["device_id"].(string),
		FirmwareVersion: pktData["firmware_version"].(float64),
		NetworkType:     networkType,
		HappenedAt:      happenedAt,
		Timestamp:       timestampInt, // Store parsed timestamp as int64.
		BeaconsAmount:   int(pktData["beacons_amount"].(float64)),
		RadarCumulative: int(pktData["radar_cumulative"].(float64)),
		IsOccupied:      pktData["is_occupied"].(float64) != 0, // Convert to boolean.
		Beacons:         beacons,                               // Attach the processed beacons.
	}

	if networkType == "NB-Iot" || networkType == "LoRa" {
		activityLog.MagnetAbsTotal = int(pktData["magnet_abs_total"].(float64))
		activityLog.PeakDistanceCm = int(pktData["peak_distance_cm"].(float64))
	}

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
		return fmt.Errorf("failed to execute bulk insert for activity logs: %w", err)
	}

	return nil
}
