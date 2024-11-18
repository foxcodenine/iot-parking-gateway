-- Create the 'parking' schema for organizing related tables and functions.
CREATE SCHEMA parking;

-- ---------------------------------------------------------------------

-- Enable TimescaleDB for time-series data management.
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- Set the transaction isolation level and timezone to ensure data consistency and uniform time representation.
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
SET TIMEZONE='UTC';

-- ---------------------------------------------------------------------

-- Devices table stores unique device identifiers and tracks their creation and update timestamps.
-- Create a table to store device information, including network details and location
CREATE TABLE parking.devices (
    device_id VARCHAR(255) PRIMARY KEY,  -- IMEI, typically 15 digits.
    name VARCHAR(255) NULL,              -- Name of the device.
    network_type VARCHAR(50) NULL,       -- Type of network (e.g., NB-IoT, LoRa).
    firmware_version DECIMAL(5, 2) NULL, -- Firmware version of the device.
    latitude DECIMAL(9, 6) DEFAULT 0,    -- Latitude for geographic location.
    longitude DECIMAL(9, 6) DEFAULT 0,   -- Longitude for geographic location.
    beacons JSONB,                  -- JSONB data type for beacon information.
    happened_at TIMESTAMP NULL,          -- Timestamp of the device data capture.
    occupied BOOLEAN DEFAULT FALSE,      -- Whether the space is occupied.
    created_at TIMESTAMP DEFAULT NOW(),  -- Set at record creation.
    updated_at TIMESTAMP DEFAULT NOW()   -- Updated automatically via trigger.
);

-- Create a function to automatically update the 'updated_at' timestamp.
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach a trigger to update the 'updated_at' field before any update operation on 'devices'.
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON parking.devices
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


-- ---------------------------------------------------------------------

-- Raw data logs table stores raw data entries with associated device info and timestamps.
CREATE TABLE parking.raw_data_logs (
    id UUID,
    device_id VARCHAR(500),                -- Optional device identifier.
    firmware_version DECIMAL(5, 2) NOT NULL,-- Required firmware version.
    network_type VARCHAR(50) NOT NULL,      -- Required network type.
    raw_data TEXT NOT NULL,                 -- Raw data text.
    created_at TIMESTAMP DEFAULT NOW(),     -- Creation timestamp.
    PRIMARY KEY (id)
);

-- ---------------------------------------------------------------------

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
    occupied BOOLEAN,
    beacons JSONB,                         -- JSONB column for beacon data.
    PRIMARY KEY (id, happened_at)
);

-- Convert 'activity_logs' to a TimescaleDB hypertable for optimized time-series queries.
SELECT create_hypertable('parking.activity_logs', 'happened_at');

-- Index the 'beacons' JSONB column to improve query performance on JSON data.
-- CREATE INDEX ON parking.activity_logs USING GIN (beacons);

-- Optional: Indexes for improving query performance on frequently queried columns
CREATE INDEX idx_activity_logs_device_id ON parking.activity_logs (device_id);
CREATE INDEX idx_activity_logs_happened_at ON parking.activity_logs (happened_at);

-- Verify the configuration of hypertables in the database.
SELECT * FROM timescaledb_information.hypertables;

-- ---------------------------------------------------------------------
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

-- Optional: Indexes for improving query performance on frequently queried columns
CREATE INDEX idx_nbiot_keepalive_logs_device_id ON parking.nbiot_keepalive_logs (device_id);
CREATE INDEX idx_nbiot_keepalive_logs_happened_at ON parking.nbiot_keepalive_logs (happened_at);

-- ---------------------------------------------------------------------

CREATE TABLE parking.nbiot_setting_logs (
    id SERIAL PRIMARY KEY,
    raw_id UUID NOT NULL,
    device_id VARCHAR(255) NOT NULL,               -- Device identifier, can be IMEI or UUID
    happened_at TIMESTAMP NOT NULL,                -- Timestamp when the settings event occurred
    created_at TIMESTAMP DEFAULT NOW(),            -- Timestamp when the record was created
    timestamp BIGINT,

    device_mode SMALLINT,
    device_enable SMALLINT,
    radar_car_cal_lo_th INTEGER,
    radar_car_cal_hi_th INTEGER,
    radar_car_uncal_lo_th INTEGER,
    radar_car_uncal_hi_th INTEGER,
    radar_car_delta_th INTEGER,
    mag_car_lo INTEGER,
    mag_car_hi INTEGER,
    radar_trail_cal_lo_th INTEGER,
    radar_trail_cal_hi_th INTEGER,
    radar_trail_uncal_lo_th INTEGER,
    radar_trail_uncal_hi_th INTEGER,
    debug_period SMALLINT,
    debug_mode SMALLINT,
    logs_mode SMALLINT,
    logs_amount SMALLINT,
    maximum_registration_time SMALLINT,
    maximum_registration_attempts SMALLINT,
    maximum_deep_sleep_time SMALLINT,
    ten_x_deep_sleep_time BIGINT,
    ten_x_action_before BIGINT,
    ten_x_action_after BIGINT,
    nb_iot_udp_ip VARCHAR(15),
    nb_iot_udp_port INTEGER,
    nb_iot_apn_length SMALLINT,
    nb_iot_apn VARCHAR(255),
    nb_iot_imsi VARCHAR(15)
);

-- Optional: Indexes for improving query performance on frequently queried columns
CREATE INDEX idx_nbiot_setting_logs_device_id ON parking.nbiot_setting_logs (device_id);
CREATE INDEX idx_nbiot_setting_logs_happened_at ON parking.nbiot_setting_logs (happened_at);
