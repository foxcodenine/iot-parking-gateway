-- NB-IoT keepalive logs table tracks keepalive events and various metrics.
CREATE TABLE IF NOT EXISTS parking.sigfox_keepalive_logs (
    id SERIAL,
    raw_id UUID NOT NULL,
    device_id VARCHAR(255) NOT NULL,               -- Required device identifier.
    firmware_version DECIMAL(5, 2) NOT NULL,       -- Required firmware version.
    network_type VARCHAR(50) NOT NULL,             -- Required network type.
    happened_at TIMESTAMP NOT NULL,                -- Timestamp of keepalive event.
    created_at TIMESTAMP DEFAULT NOW(),            -- Timestamp when the record is created.
    timestamp BIGINT,                              -- Event timestamp in UNIX format.

    idle_voltage SMALLINT,                         -- Idle voltage in V.
    battery_percentage SMALLINT DEFAULT NULL,               
    current SMALLINT,                             -- Current in mA.
    reset_count SMALLINT,                         -- Reset count.    
    temperature_min SMALLINT,                     -- Minimum temperature.
    temperature_max SMALLINT,                     -- Maximum temperature.
    radar_error SMALLINT,                         -- Radar error count.   
    tcve_error SMALLINT,                          -- TCVE error count.
    radar_cumulative INT,                       
    settings_checksum SMALLINT DEFAULT NULL,                   -- Settings checksum value.

    PRIMARY KEY (id, happened_at)
);

-- Convert 'sigfox_keepalive_logs' to a TimescaleDB hypertable for optimized performance.
SELECT create_hypertable('parking.sigfox_keepalive_logs', 'happened_at');

-- Indexes for improved query performance.
CREATE INDEX idx_sigfox_keepalive_logs_device_id ON parking.sigfox_keepalive_logs (device_id);
CREATE INDEX idx_sigfox_keepalive_logs_happened_at ON parking.sigfox_keepalive_logs (happened_at);
