package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	up "github.com/upper/db/v4"
)

// -----------------------------------------------------------------------------

// Device represents a parking device in the database with ID and timestamps.
type Device struct {
	DeviceID  string    `db:"device_id" json:"device_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// -----------------------------------------------------------------------------

// TableName returns the full table name for the Device model in PostgreSQL.
func (d *Device) TableName() string {
	return "parking.devices"
}

// -----------------------------------------------------------------------------

// GetAll retrieves all devices from the database.
func (d *Device) GetAll() ([]Device, error) {
	collection := dbSession.Collection(d.TableName())

	var devices []Device
	err := collection.Find().All(&devices)
	if err != nil {
		return nil, err
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
func (d *Device) Create(id string) (*Device, error) {
	collection := dbSession.Collection(d.TableName())

	// Set current time for CreatedAt and UpdatedAt
	now := time.Now().UTC()

	// Prepare the new device record
	newDevice := Device{
		DeviceID:  id,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := collection.InsertReturning(&newDevice)
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

	return &newDevice, nil
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
