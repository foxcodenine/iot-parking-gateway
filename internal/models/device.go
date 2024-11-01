package models

import (
	"errors"
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

	newDevice := Device{
		DeviceID:  id,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err := collection.Insert(newDevice)
	if err != nil {
		return nil, err
	}

	return &newDevice, nil
}

// -----------------------------------------------------------------------------

// UpdateByID updates an existing device by its ID.
func (d *Device) UpdateByID(id string, updatedFields map[string]interface{}) (*Device, error) {
	collection := dbSession.Collection(d.TableName())

	updatedFields["updated_at"] = time.Now().UTC()

	err := collection.Find(up.Cond{"device_id": id}).Update(updatedFields)
	if err != nil {
		return nil, err
	}

	return d.GetByID(id) // return the updated device
}

// DeleteByID deletes a device by its ID.
func (d *Device) DeleteByID(id string) error {
	collection := dbSession.Collection(d.TableName())

	err := collection.Find(up.Cond{"device_id": id}).Delete()
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
