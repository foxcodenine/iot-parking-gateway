-- NB-IoT setting logs table tracks configuration changes.
CREATE TABLE parking.sigfox_setting_logs (
    id SERIAL,
    raw_id UUID NOT NULL,
    device_id VARCHAR(255) NOT NULL,               -- Device identifier, can be IMEI or UUID
    firmware_version DECIMAL(5, 2) NOT NULL,       -- Required firmware version.
    network_type VARCHAR(50) NOT NULL,             -- Required network type.
    happened_at TIMESTAMP NOT NULL,                -- Timestamp when the settings event occurred
    created_at TIMESTAMP DEFAULT NOW(),            -- Timestamp when the record was created
    timestamp BIGINT,

    device_mode SMALLINT,
    device_enable SMALLINT,
    radar_car_cal_lo_th INTEGER,
    radar_car_cal_hi_th INTEGER,
    radar_car_delta_th INTEGER,
    downlink_en_7_bits_repeated_occupancy_period_mins SMALLINT,

    PRIMARY KEY (id, happened_at)
);
-- Convert 'sigfox_setting_logs' to a TimescaleDB hypertable for optimized performance.
SELECT create_hypertable('parking.sigfox_setting_logs', 'happened_at');

-- Indexes for improved query performance.
CREATE INDEX idx_sigfox_setting_logs_device_id ON parking.sigfox_setting_logs (device_id);
CREATE INDEX idx_sigfox_setting_logs_happened_at ON parking.sigfox_setting_logs (happened_at);
