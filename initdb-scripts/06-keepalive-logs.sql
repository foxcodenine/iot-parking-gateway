-- NB-IoT keepalive logs table tracks keepalive events and various metrics.
CREATE TABLE IF NOT EXISTS parking.nbiot_keepalive_logs (
    id SERIAL,
    raw_id UUID NOT NULL,
    device_id VARCHAR(255) NOT NULL,               -- Required device identifier.
    firmware_version DECIMAL(5, 2) NOT NULL,       -- Required firmware version.
    network_type VARCHAR(50) NOT NULL,             -- Required network type.
    happened_at TIMESTAMP NOT NULL,                -- Timestamp of keepalive event.
    created_at TIMESTAMP DEFAULT NOW(),            -- Timestamp when the record is created.
    timestamp BIGINT,                              -- Event timestamp in UNIX format.

    idle_voltage SMALLINT,                         -- Idle voltage in V.
    battery_percentage SMALLINT,                  -- Battery percentage.
    current SMALLINT,                             -- Current in mA.
    reset_count SMALLINT,                         -- Reset count.
    manual_calibration BOOLEAN,                   -- Manual calibration status.
    temperature_min SMALLINT,                     -- Minimum temperature.
    temperature_max SMALLINT,                     -- Maximum temperature.
    radar_error SMALLINT,                         -- Radar error count.
    mag_error SMALLINT,                           -- Magnetometer error count.
    tcve_error SMALLINT,                          -- TCVE error count.
    ble_security_issues SMALLINT,                 -- BLE security issues count.
    radar_cumulative_total INTEGER,               -- Radar cumulative total.
    mag_total INTEGER,                            -- Magnetometer total value.
    network_registration_ok SMALLINT,             -- Registration successful count.
    network_registration_nok SMALLINT,            -- Registration failed count.
    rssi_average SMALLINT,                        -- RSSI average value.
    network_message_attempts SMALLINT,            -- Message attempts count.
    network_ack_1ds SMALLINT,                     -- Network ACK 1DS count.
    network_1ack_ds SMALLINT,                     -- Network 1ACK DS count.
    network_1ack_1ds SMALLINT,                    -- Network 1ACK 1DS count.
    tcvr_deep_sleep_min SMALLINT,                 -- Deep sleep min time.
    tcvr_deep_sleep_max SMALLINT,                 -- Deep sleep max time.
    tcvr_deep_sleep_average SMALLINT,             -- Deep sleep average time.
    settings_checksum SMALLINT,                   -- Settings checksum value.
    socket_error SMALLINT,                        -- Socket error count.
    t3324 SMALLINT,                               -- Timer T3324 value.
    t3412 SMALLINT,                               -- Timer T3412 value.
    time_sync_rand_byte SMALLINT,                 -- Time sync random byte.
    time_sync_current_unix_time BIGINT,           -- Time sync random byte.


    PRIMARY KEY (id, happened_at)
);

-- Convert 'nbiot_keepalive_logs' to a TimescaleDB hypertable for optimized performance.
SELECT create_hypertable('parking.nbiot_keepalive_logs', 'happened_at');

-- Indexes for improved query performance.
CREATE INDEX idx_nbiot_keepalive_logs_device_id ON parking.nbiot_keepalive_logs (device_id);
CREATE INDEX idx_nbiot_keepalive_logs_happened_at ON parking.nbiot_keepalive_logs (happened_at);
