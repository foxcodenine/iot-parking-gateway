package models

import "time"

// NbDeviceData represents the structured data for device updates.
type NbDeviceData struct {
	DeviceID        string    `json:"device_id"`
	FirmwareVersion float64   `json:"firmware_version"`
	Beacons         []Beacon  `json:"beacons"`
	HappenedAt      time.Time `json:"happened_at"`
	Occupied        bool      `json:"occupied"`
}
