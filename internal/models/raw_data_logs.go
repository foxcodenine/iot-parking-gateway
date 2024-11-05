package models

import (
	"time"

	"github.com/google/uuid"
)

type RawDataLog struct {
	Uuid            uuid.UUID `db:"uuid" json:"uuid"`
	DeviceID        string    `db:"device_id" json:"device_id"`
	FirmwareVersion int       `db:"firmware_version" json:"firmware_version"`
	NetworkType     string    `db:"network_type" json:"network_type"`
	RawData         string    `db:"raw_data" json:"raw_data"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

func (r *RawDataLog) TableName() string {
	return "parking.raw_data_logs"
}
