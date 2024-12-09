package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
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

	// Attempt to retrieve cached devices
	cachedData, err := cache.AppCache.Get("db:devices")
	if err != nil {
		helpers.LogError(err, "Failed to get devices from cache")
		return nil, err
	}

	if cachedData != nil {
		// Asserting the type of cached data to []interface{}
		cachedDevices, ok := cachedData.([]interface{})
		if !ok {
			helpers.LogError(fmt.Errorf("cache data type mismatch: expected []interface{}, got %T", cachedData), "Cache data type mismatch")
			return nil, fmt.Errorf("cache data type mismatch: expected []interface{}, got %T", cachedData)
		}

		// Initialize slice to hold the converted device objects
		devices = make([]*Device, len(cachedDevices))
		for i, cachedDevice := range cachedDevices {
			deviceMap, ok := cachedDevice.(map[string]interface{})
			if !ok {
				helpers.LogError(fmt.Errorf("failed to assert type for device data: %T", cachedDevice), "Error asserting type for cached device data")
				continue
			}

			device := &Device{} // Create a new Device instance

			// Map data from deviceMap to device struct fields
			if deviceID, ok := deviceMap["device_id"].(string); ok {
				device.DeviceID = deviceID
			}
			if name, ok := deviceMap["name"].(string); ok {
				device.Name = name
			}
			if networkType, ok := deviceMap["network_type"].(string); ok {
				device.NetworkType = networkType
			}
			if firmwareVersion, ok := deviceMap["firmware_version"].(float64); ok {
				device.FirmwareVersion = firmwareVersion
			}
			if latitude, ok := deviceMap["latitude"].(float64); ok {
				device.Latitude = latitude
			}
			if longitude, ok := deviceMap["longitude"].(float64); ok {
				device.Longitude = longitude
			}
			if beacons, ok := deviceMap["beacons"].(string); ok {
				device.BeaconsJSON = sql.NullString{String: beacons, Valid: true}
				_ = device.ParseBeaconsJSON()
			}
			if happenedAt, ok := deviceMap["happened_at"].(string); ok {
				device.HappenedAt, _ = time.Parse(time.RFC3339, happenedAt)
			}
			if isOccupied, ok := deviceMap["is_occupied"].(bool); ok {
				device.IsOccupied = isOccupied
			}
			if isAllowed, ok := deviceMap["is_allowed"].(bool); ok {
				device.IsAllowed = isAllowed
			}
			if isBlocked, ok := deviceMap["is_blocked"].(bool); ok {
				device.IsBlocked = isBlocked
			}
			if isHidden, ok := deviceMap["is_hidden"].(bool); ok {
				device.IsHidden = isHidden
			}
			if createdAt, ok := deviceMap["created_at"].(string); ok {
				device.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
			}
			if updatedAt, ok := deviceMap["updated_at"].(string); ok {
				device.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
			}

			devices[i] = device
		}

		return devices, nil
	}

	// If not cached, fetch from the database
	collection := dbSession.Collection(d.TableName())
	err = collection.Find(up.Cond{"deleted_at": nil}).OrderBy("created_at").All(&devices)

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

	// Cache the devices after successful database fetch
	ttl, err := strconv.Atoi(os.Getenv("REDIS_DEFAULT_TTL"))
	if err != nil {
		helpers.LogError(err, "Failed to convert REDIS_DEFAULT_TTL to integer")
		ttl = 600 // Default TTL as a fallback
	}

	err = cache.AppCache.Set("db:devices", devices, ttl)
	if err != nil {
		helpers.LogError(err, "Failed to set devices in cache")
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

	err = cache.AppCache.Delete("db:devices")

	if err != nil {
		helpers.LogError(err, "failed to delete devices from cache")
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
			happened_at, is_occupied, is_allowed, is_blocked, created_at, updated_at, deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NULL
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
		time.Now().UTC(),
		time.Now().UTC(),
	}

	// Execute the query
	_, err := collection.Session().SQL().Exec(sqlQuery, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert device: %w", err)
	}

	// Return the updated device
	return d.GetByID(device.DeviceID)
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

	// Clear device in db cache
	err = cache.AppCache.Delete("db:devices")
	if err != nil {
		helpers.LogError(err, "failed to delete devices from cache")
	}

	err = res.Update(updatedFields)
	if err != nil {
		return nil, fmt.Errorf("error updating device: %w", err)
	}

	return d.GetByID(id) // return the updated device
}

// DeleteByID deletes a device by its ID.
func (d *Device) DeleteByID(id string) error {
	collection := dbSession.Collection(d.TableName())

	err := cache.AppCache.Delete("db:devices")
	if err != nil {
		helpers.LogError(err, "failed to delete devices from cache")
	}

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

	// Clear device from the cache
	err = cache.AppCache.Delete("db:devices")
	if err != nil {
		helpers.LogError(err, "failed to delete devices from cache")
	}

	// Perform the soft deletion
	err = res.Update(updatedFields)
	if err != nil {
		return fmt.Errorf("error soft-deleting device: %w", err)
	}

	return nil // Successful soft deletion
}

// -----------------------------------------------------------------------------

// BulkUpdateDevices updates multiple device records based on their device IDs.
func (d *Device) BulkUpdateDevices(deviceData []ActivityLog) error {

	if len(deviceData) == 0 {
		return nil // No data to update
	}

	err := cache.AppCache.Delete("db:devices")
	if err != nil {
		helpers.LogError(err, "failed to delete devices from cache")
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
	_, err = dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk update for devices: %w", err)
	}

	return nil
}
