package models

import (
	"fmt"
	"time"
	// up "github.com/upper/db/v4"
)

type Device struct {
	DeviceID  string    `db:"device_id" json:"device_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the PostgreSQL table name for the Device model.
func (d *Device) TableName() string {

	return "parking.devices"
}

func (d *Device) Create(id string) (*Device, error) {
	collection := upper.Collection(d.TableName())

	newDevice := Device{
		DeviceID: id,
	}

	res, err := collection.Insert(newDevice)
	if err != nil {
		return nil, err
	}

	fmt.Println(res)

	return &newDevice, nil

}
