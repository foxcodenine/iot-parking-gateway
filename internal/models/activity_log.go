package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID               uuid.UUID `db:"id" json:"id"`
	RawID            uuid.UUID `db:"raw_id" json:"raw_id"`
	HappenedAt       time.Time `db:"happened_at" json:"happened_at"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	Timestamp        int64     `db:"timestamp" json:"timestamp"`
	BeaconsAmount    int       `db:"beacons_amount" json:"beacons_amount"`
	MagnetAbsTotal   int       `db:"magnet_abs_total" json:"magnet_abs_total"`
	NextOffset       int       `db:"next_offset" json:"next_offset"`
	PeakDistanceCm   int       `db:"peak_distance_cm" json:"peak_distance_cm"`
	RadarCumulative  int       `db:"radar_cumulative" json:"radar_cumulative"`
	VehicleOccupancy bool      `db:"vehicle_occupancy" json:"vehicle_occupancy"`
}

func (a *ActivityLog) TableName() string {
	return "parking.activity_logs"
}
