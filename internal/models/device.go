package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	up "github.com/upper/db/v4"
)

// -----------------------------------------------------------------------------

// Device represents a parking device in the database with ID and timestamps.
type Device struct {
	DeviceID        string         `db:"device_id" json:"device_id"`
	Name            string         `db:"name" json:"name"`
	NetworkType     string         `db:"network_type,omitempty" json:"network_type"`
	FirmwareVersion float64        `db:"firmware_version" json:"firmware_version"`
	Latitude        float64        `db:"latitude" json:"latitude"`
	Longitude       float64        `db:"longitude" json:"longitude"`
	BeaconsJSON     sql.NullString `db:"beacons,omitempty" json:"beacons_json"`
	Beacons         []Beacon       `json:"beacons"`
	HappenedAt      time.Time      `db:"happened_at" json:"happened_at"`
	IsOccupied      bool           `db:"is_occupied" json:"is_occupied"`
	IsAllowed       bool           `db:"is_allowed" json:"is_allowed"` // Indicates if the device is allowed
	IsBlocked       bool           `db:"is_blocked" json:"is_blocked"` // Indicates if the device is blocked
	IsHidden        bool           `db:"is_hidden" json:"is_hidden"`   // Indicates if the device is hidden
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at"`
}

// -----------------------------------------------------------------------------

// TableName returns the full table name for the Device model in PostgreSQL.
func (d *Device) TableName() string {
	return "parking.devices"
}

// -----------------------------------------------------------------------------

// ParseBeaconsJSON parses the BeaconsJSON field into the Beacons field.
func (d *Device) ParseBeaconsJSON() error {

	if !d.BeaconsJSON.Valid {
		// If BeaconsJSON is NULL, skip parsing
		d.Beacons = nil
		return nil
	}

	var decodedBeacons []Beacon
	err := json.Unmarshal([]byte(d.BeaconsJSON.String), &decodedBeacons)
	if err != nil {
		helpers.LogError(err, "Error unmarshalling JSON:")
		return err
	}
	d.Beacons = decodedBeacons

	return nil
}

// GetAll retrieves all devices from the database.
func (d *Device) GetAll() ([]Device, error) {
	collection := dbSession.Collection(d.TableName())

	var devices []Device
	err := collection.Find().All(&devices)
	if err != nil {
		return nil, err
	}

	for i := range devices {
		err := devices[i].ParseBeaconsJSON()
		if err != nil {
			return nil, fmt.Errorf("error parsing BeaconsJSON for device ID %s: %w", devices[i].DeviceID, err)
		}
	}

	return devices, nil
}

// -----------------------------------------------------------------------------

// GetByID retrieves a single device by its ID.
func (d *Device) GetByID(id string) (*Device, error) {
	collection := dbSession.Collection(d.TableName())

	var device Device

	err := collection.Find(up.Cond{"device_id": id}).One(&device)
	if err != nil {
		if errors.Is(err, up.ErrNoMoreRows) {
			return nil, errors.New("device not found")
		}
		return nil, err
	}

	return &device, nil
}

// -----------------------------------------------------------------------------

// Create inserts a new Device record in the database and returns the created device.
func (d *Device) Create(newDevice *Device) (*Device, error) {
	collection := dbSession.Collection(d.TableName())

	// Set current time for CreatedAt and UpdatedAt
	now := time.Now().UTC()
	newDevice.CreatedAt = now
	newDevice.UpdatedAt = now

	_, err := collection.Insert(newDevice)
	if err != nil {
		// Check if the error is a duplicate key violation (PostgreSQL SQL state 23505)
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			// Define a custom error for duplicate device entries
			var ErrDuplicateDevice = errors.New("device with this ID already exists")
			return nil, ErrDuplicateDevice
		}
		// Return any other errors with additional context
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	return newDevice, nil
}

// -----------------------------------------------------------------------------

// UpdateByID updates an existing device by its ID.
func (d *Device) UpdateByID(id string, updatedFields map[string]interface{}) (*Device, error) {
	collection := dbSession.Collection(d.TableName())

	updatedFields["updated_at"] = time.Now().UTC()

	// Check if the device exists by counting matching rows
	res := collection.Find(up.Cond{"device_id": id})
	count, err := res.Count()
	if err != nil {
		return nil, fmt.Errorf("error checking device existence: %w", err)
	}

	// If no matching device is found, return a custom error
	if count == 0 {
		return nil, fmt.Errorf("device with ID %s not found", id)
	}

	err = res.Update(updatedFields)
	if err != nil {
		return nil, fmt.Errorf("error updating device: %w", err)
	}

	// Check if `device_id` was updated; if not, return an error
	newID, ok := updatedFields["device_id"].(string)
	if !ok {
		return nil, fmt.Errorf("updated device with ID: %s not found", newID)
	}

	return d.GetByID(newID) // return the updated device
}

// DeleteByID deletes a device by its ID.
func (d *Device) DeleteByID(id string) error {
	collection := dbSession.Collection(d.TableName())

	// Check if the device exists by counting matching rows
	res := collection.Find(up.Cond{"device_id": id})
	count, err := res.Count()
	if err != nil {
		return fmt.Errorf("error checking device existence: %w", err)
	}

	// If no matching device is found, return a custom error
	if count == 0 {
		return errors.New("device not found")
	}

	// Proceed with deletion since the device exists
	if err := res.Delete(); err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	return nil // Successful deletion
}

// -----------------------------------------------------------------------------

// BulkUpdateDevices updates multiple device records based on their device IDs.
func (d *Device) BulkUpdateDevices(deviceData []ActivityLog) error {

	if len(deviceData) == 0 {
		return nil // No data to update
	}

	var args []interface{}
	valuesList := make([]string, len(deviceData))
	for i, data := range deviceData {
		// Prepare a position offset for SQL placeholders based on number of fields per device
		pos := i*5 + 1
		valuesList[i] = fmt.Sprintf("($%d, $%d::numeric, $%d::jsonb, $%d::timestamp, $%d::boolean)", pos, pos+1, pos+2, pos+3, pos+4)
		args = append(args, data.DeviceID, data.FirmwareVersion, data.Beacons, data.HappenedAt, data.IsOccupied)
	}

	// Construct the SQL statement with explicit type casts to ensure proper data handling
	query := fmt.Sprintf(`
		UPDATE parking.devices AS d
		SET
			firmware_version = v.firmware_version,
			beacons = v.beacons,
			happened_at = v.happened_at,
			is_occupied = v.is_occupied,
			updated_at = NOW()
		FROM (VALUES
			%s
		) AS v(device_id, firmware_version, beacons, happened_at, is_occupied)
		WHERE d.device_id = v.device_id
	`, strings.Join(valuesList, ", "))

	// Execute the constructed query with the arguments
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk update for devices: %w", err)
	}

	return nil
}
