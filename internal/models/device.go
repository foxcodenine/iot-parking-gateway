package models

import (
	"time"
	// up "github.com/upper/db/v4"
)

// Device represents a parking device in the database with ID and timestamps.
type Device struct {
	DeviceID  string    `db:"device_id" json:"device_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the full table name for the Device model in PostgreSQL.
func (d *Device) TableName() string {
	return "parking.devices"
}

// Create inserts a new Device record in the database and returns the created device.
func (d *Device) Create(id string) (*Device, error) {
	collection := upper.Collection(d.TableName())

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
