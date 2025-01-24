package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
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
	KeepaliveAt     time.Time      `db:"keepalive_at" json:"keepalive_at"`
	SettingsAt      time.Time      `db:"settings_at" json:"settings_at"`
	IsOccupied      bool           `db:"is_occupied" json:"is_occupied"`
	IsAllowed       bool           `db:"is_allowed" json:"is_allowed"` // Indicates if the device is allowed
	IsBlocked       bool           `db:"is_blocked" json:"is_blocked"` // Indicates if the device is blocked
	IsHidden        bool           `db:"is_hidden" json:"is_hidden"`   // Indicates if the device is hidden
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at"`
	DeletedAt       time.Time      `db:"deleted_at,omitempty" json:"deleted_at,omitempty"`
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
func (d *Device) GetAll() ([]*Device, error) {
	var devices []*Device

	collection := dbSession.Collection(d.TableName())
	err := collection.Find(up.Cond{"deleted_at": nil}).OrderBy("created_at").All(&devices)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve devices from database: %w", err)
	}

	for i := range devices {
		err := devices[i].ParseBeaconsJSON()
		if err != nil {
			helpers.LogError(err, fmt.Sprintf("Error parsing BeaconsJSON for device ID %s:", devices[i].DeviceID))
			return nil, fmt.Errorf("error parsing BeaconsJSON for device ID %s: %w", devices[i].DeviceID, err)
		}
	}

	return devices, nil
}

// -----------------------------------------------------------------------------

// GetAllIncludingDeleted retrieves all devices from the database, including those that have been soft-deleted.
func (d *Device) GetAllIncludingDeleted() ([]*Device, error) {
	var devices []*Device

	// Fetch all devices, including those with `deleted_at` not NULL
	collection := dbSession.Collection(d.TableName())
	err := collection.Find().OrderBy("created_at").All(&devices)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all devices including deleted ones: %w", err)
	}

	// Parse BeaconsJSON for each device
	for i := range devices {
		err := devices[i].ParseBeaconsJSON()
		if err != nil {
			helpers.LogError(err, fmt.Sprintf("Error parsing BeaconsJSON for device ID %s:", devices[i].DeviceID))
			return nil, fmt.Errorf("error parsing BeaconsJSON for device ID %s: %w", devices[i].DeviceID, err)
		}
	}

	return devices, nil
}

// -----------------------------------------------------------------------------

// GetByID retrieves a single device by its ID, excluding soft-deleted records.
func (d *Device) GetByID(id string) (*Device, error) {
	collection := dbSession.Collection(d.TableName())

	var device Device

	// Add condition to check `deleted_at IS NULL` in addition to `device_id`.
	err := collection.Find(up.Cond{"device_id": id, "deleted_at": nil}).One(&device)
	if err != nil {
		if errors.Is(err, up.ErrNoMoreRows) {
			return nil, errors.New("device not found or has been deleted")
		}
		return nil, fmt.Errorf("failed to retrieve device: %w", err)
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

	// Convert the newDevice struct to a map
	deviceData, err := helpers.StructToMap(newDevice)
	if err != nil {
		return nil, fmt.Errorf("failed to convert device data: %w", err)
	}

	if _, ok := deviceData["keepalive_at"]; !ok {

		deviceData["keepalive_at"] = "0001-01-01T00:00:00Z"
	}

	if _, ok := deviceData["settings_at"]; !ok {

		deviceData["settings_at"] = "0001-01-01T00:00:00Z"
	}
	if _, ok := deviceData["happened_at"]; !ok {

		deviceData["happened_at"] = "0001-01-01T00:00:00Z"
	}

	// Save the device data to the cache
	if err := cache.AppCache.SaveDeviceData(newDevice.DeviceID, deviceData); err != nil {
		return nil, fmt.Errorf("failed to cache device data: %w", err)
	}

	return newDevice, nil
}

// Upsert inserts a new device or updates specific fields if the device already exists.
func (d *Device) Upsert(device *Device) (*Device, error) {
	collection := dbSession.Collection(d.TableName())

	// Prepare the upsert query with positional placeholders
	sqlQuery := `
		INSERT INTO parking.devices (
			device_id, name, network_type, firmware_version, latitude, longitude, beacons, 
			happened_at, is_occupied, is_allowed, is_blocked, is_hidden, created_at, updated_at, deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,  $14, NULL
		)
		ON CONFLICT (device_id) DO UPDATE SET
			is_allowed = false,
			is_blocked = false,
			is_hidden = false,
			deleted_at = NULL,
			updated_at = EXCLUDED.updated_at;
	`

	// Prepare the parameter values
	params := []interface{}{
		device.DeviceID,
		device.Name,
		device.NetworkType,
		device.FirmwareVersion,
		device.Latitude,
		device.Longitude,
		device.BeaconsJSON,
		device.HappenedAt,
		device.IsOccupied,
		device.IsAllowed,
		device.IsBlocked,
		device.IsHidden,
		time.Now().UTC(),
		time.Now().UTC(),
	}

	// Execute the query
	_, err := collection.Session().SQL().Exec(sqlQuery, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert device: %w", err)
	}

	// Fetch the updated device from the database
	upsertDevice, err := d.GetByID(device.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated device: %w", err)
	}

	// Convert the updated device struct to a map
	deviceData, err := helpers.StructToMap(upsertDevice)
	if err != nil {
		return nil, fmt.Errorf("failed to convert device data: %w", err)
	}

	if _, ok := deviceData["keepalive_at"]; !ok {

		deviceData["keepalive_at"] = "0001-01-01T00:00:00Z"
	}

	if _, ok := deviceData["settings_at"]; !ok {

		deviceData["settings_at"] = "0001-01-01T00:00:00Z"
	}
	if _, ok := deviceData["happened_at"]; !ok {

		deviceData["happened_at"] = "0001-01-01T00:00:00Z"
	}

	// Save the device data to the cache
	if err := cache.AppCache.SaveDeviceData(upsertDevice.DeviceID, deviceData); err != nil {
		return nil, fmt.Errorf("failed to cache device data: %w", err)
	}

	return upsertDevice, nil

}

// -----------------------------------------------------------------------------

// UpdateByID updates an existing device by its ID.
func (d *Device) UpdateByID(id string, updatedFields map[string]interface{}) (*Device, error) {
	collection := dbSession.Collection(d.TableName())

	updatedFields["updated_at"] = time.Now().UTC()

	// // Explicitly handle `deleted_at` to ensure it is set to NULL
	// updatedFields["deleted_at"] = up.Raw("NULL")

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

	// Update the device data to the cache
	if err := cache.AppCache.UpdateDeviceFields(id, updatedFields); err != nil {
		return nil, fmt.Errorf("failed to update device data in cache: %w", err)
	}

	return d.GetByID(id) // return the updated device
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

	// Delete the device data from cache
	if err := cache.AppCache.DeleteDevice(id); err != nil {
		return fmt.Errorf("failed to delete device data from cache: %w", err)
	}

	return nil // Successful deletion
}

// SoftDeleteByID marks a device as deleted by setting the `deleted_at` field to the current timestamp.
func (d *Device) SoftDeleteByID(id string) error {
	collection := dbSession.Collection(d.TableName())

	// Check if the device exists
	res := collection.Find(up.Cond{"device_id": id})
	count, err := res.Count()
	if err != nil {
		return fmt.Errorf("error checking device existence: %w", err)
	}

	// If no matching device is found, return a custom error
	if count == 0 {
		return fmt.Errorf("device with ID %s not found", id)
	}

	// Prepare the updated fields for soft deletion
	updatedFields := map[string]interface{}{
		"deleted_at": time.Now().UTC(),
		"updated_at": time.Now().UTC(),
		"is_blocked": true,
	}

	// Perform the soft deletion
	err = res.Update(updatedFields)
	if err != nil {
		return fmt.Errorf("error soft-deleting device: %w", err)
	}

	if err := cache.AppCache.UpdateDeviceFields(id, updatedFields); err != nil {
		return fmt.Errorf("failed to soft delete device data in cache: %w", err)
	}

	return nil // Successful soft deletion
}

// -----------------------------------------------------------------------------

// BulkUpdateDevices updates multiple device records based on their device IDs.
func (d *Device) BulkUpdateDevices(deviceData []Device) error {

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
			deleted_at = NULL,
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

// BulkUpdateDevicesKeepalive updates the `keepalive_at` field for multiple devices in a single database query.
// This method ensures efficient updates by batching updates into a single SQL statement.
func (d *Device) BulkUpdateDevicesKeepalive(deviceData []Device) error {

	// Return early if no data is provided to avoid unnecessary processing.
	if len(deviceData) == 0 {
		return nil // No data to update
	}

	var args []interface{}                        // Arguments to be passed to the query.
	valuesList := make([]string, len(deviceData)) // List of value placeholders for SQL.

	for i, data := range deviceData {
		// Prepare a position offset for SQL placeholders based on the number of fields per device.
		pos := i*2 + 1 // Each device requires 2 placeholders: `device_id` and `keepalive_at`.
		valuesList[i] = fmt.Sprintf("($%d, $%d::timestamp)", pos, pos+1)
		args = append(args, data.DeviceID, data.KeepaliveAt) // Add the actual values.
	}

	// Construct the SQL statement with explicit type casts to ensure proper data handling.
	query := fmt.Sprintf(`
		UPDATE parking.devices AS d
		SET
			keepalive_at = v.keepalive_at,
			updated_at = NOW() 
		FROM (VALUES
			%s
		) AS v(device_id, keepalive_at)
		WHERE d.device_id = v.device_id
	`, strings.Join(valuesList, ", "))

	// Execute the constructed query with the arguments.
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk update for keepalive_at field: %w", err)
	}

	return nil // Return nil if the operation was successful.
}

// BulkUpdateDevicesSettings updates the `settings_at` field for multiple devices in a single database query.
// This method ensures efficient updates by batching updates into a single SQL statement.
func (d *Device) BulkUpdateDevicesSettings(deviceData []Device) error {
	// Return early if no data is provided to avoid unnecessary processing.
	if len(deviceData) == 0 {
		return nil // No data to update
	}

	var args []interface{}                        // Arguments to be passed to the query.
	valuesList := make([]string, len(deviceData)) // List of value placeholders for SQL.

	for i, data := range deviceData {
		// Prepare a position offset for SQL placeholders based on the number of fields per device.
		pos := i*2 + 1 // Each device requires 2 placeholders: `device_id` and `settings_at`.
		valuesList[i] = fmt.Sprintf("($%d, $%d::timestamp)", pos, pos+1)
		args = append(args, data.DeviceID, data.SettingsAt) // Add the actual values.
	}

	// Construct the SQL statement with explicit type casts to ensure proper data handling.
	query := fmt.Sprintf(`
		UPDATE parking.devices AS d
		SET
			settings_at = v.settings_at,
			updated_at = NOW() 
		FROM (VALUES
			%s
		) AS v(device_id, settings_at)
		WHERE d.device_id = v.device_id
	`, strings.Join(valuesList, ", "))

	// Execute the constructed query with the arguments.
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk update for settings_at field: %w", err)
	}

	return nil // Return nil if the operation was successful.
}
