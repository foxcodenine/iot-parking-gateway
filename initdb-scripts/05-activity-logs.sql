-- Activity logs table tracks detailed device activities and environmental metrics.
CREATE TABLE IF NOT EXISTS parking.activity_logs (
    id SERIAL,
    raw_id UUID NOT NULL,
    device_id VARCHAR(255) NOT NULL,       -- Required device identifier.
    firmware_version DECIMAL(5, 2) NOT NULL,-- Required firmware version.
    network_type VARCHAR(50) NOT NULL,     -- Required network type.
    happened_at TIMESTAMP NOT NULL,        -- Timestamp of activity occurrence.
    created_at TIMESTAMP DEFAULT NOW(),
    timestamp BIGINT,
    beacons_amount INTEGER NOT NULL DEFAULT 0,
    magnet_abs_total INTEGER,
    peak_distance_cm INTEGER,
    radar_cumulative INTEGER,
    is_occupied BOOLEAN,
    beacons JSONB,                         -- JSONB column for beacon data.
    PRIMARY KEY (id, happened_at)
);

-- Convert 'activity_logs' to a TimescaleDB hypertable for optimized time-series queries.
SELECT create_hypertable('parking.activity_logs', 'happened_at');

-- Indexes for improved query performance.
CREATE INDEX idx_activity_logs_device_id ON parking.activity_logs (device_id);
CREATE INDEX idx_activity_logs_happened_at ON parking.activity_logs (happened_at);

-- Enable Compression on the Hypertable:
ALTER TABLE parking.activity_logs SET (timescaledb.compress = true);

-- (Optional) Configure Compression Settings:
ALTER TABLE parking.activity_logs SET (timescaledb.compress_orderby = 'happened_at');
ALTER TABLE parking.activity_logs SET (timescaledb.compress_segmentby = 'device_id');

-- Add a Compression Policy:
SELECT add_compression_policy('parking.activity_logs', INTERVAL '2 months');

-- Add a Retention Policy:
SELECT add_retention_policy('parking.activity_logs', INTERVAL '12 months');

-- Remove the Retention Policy:
-- SELECT remove_retention_policy('parking.activity_logs');

