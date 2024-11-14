package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	up "github.com/upper/db/v4"
)

// Beacon represents a beacon entry associated with an activity log.
type BeaconLog struct {
	ActivityID   uuid.UUID `db:"activity_id" json:"activity_id"`     // Foreign key linking to an ActivityLog
	HappenedAt   time.Time `db:"happened_at" json:"happened_at"`     // Timestamp matching the activity's timestamp for alignment
	BeaconNumber int       `db:"beacon_number" json:"beacon_number"` // Unique number for each beacon within an activity
	Major        int       `db:"major" json:"major"`                 // Major identifier of the beacon
	Minor        int       `db:"minor" json:"minor"`                 // Minor identifier of the beacon
	RSSI         int       `db:"rssi" json:"rssi"`                   // Received Signal Strength Indicator (RSSI) value
}

// TableName returns the table name for the Beacon model.
func (b *BeaconLog) TableName() string {
	return "parking.beacon_logs"
}

// BulkInsert inserts multiple Beacon records in a single operation.
func (b *BeaconLog) BulkInsert(beacons []BeaconLog) error {
	// If there are no records to insert, exit early
	if len(beacons) == 0 {
		return nil
	}

	// Prepare slices for SQL values and arguments.
	values := make([]string, 0, len(beacons))      // Holds the placeholder for each row
	args := make([]interface{}, 0, len(beacons)*6) // Holds the actual values for each column

	for i, beacon := range beacons {
		// For each beacon, create a placeholder with indexed arguments, e.g., ($1, $2, $3, $4, $5, $6)
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))

		// Append the actual values for each placeholder in the same order as the columns
		args = append(args, beacon.ActivityID, beacon.HappenedAt, beacon.BeaconNumber, beacon.Major, beacon.Minor, beacon.RSSI)
	}

	// Construct the SQL statement by joining the placeholders for each record
	query := fmt.Sprintf("INSERT INTO %s (activity_id, happened_at, beacon_number, major, minor, rssi) VALUES %s",
		b.TableName(), strings.Join(values, ", "))

	// Execute the constructed query with the arguments
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk insert for beacons: %w", err)
	}

	return nil
}

// GetByID retrieves a single device by its ID.
func (d *Device) GetByIDs(id string) ([]BeaconLog, error) {
	collection := dbSession.Collection(d.TableName())

	var beacons []BeaconLog

	err := collection.Find(up.Cond{"activity_id": id}).All(&beacons)
	if err != nil {
		if errors.Is(err, up.ErrNoMoreRows) {
			return nil, errors.New("device not found")
		}
		return nil, err
	}

	return beacons, nil
}
