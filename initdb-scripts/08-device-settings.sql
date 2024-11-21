-- NB-IoT device settings table stores device configurations.
CREATE TABLE parking.nbiot_device_settings (
    device_id VARCHAR(255) NOT NULL PRIMARY KEY, -- Device identifier, unique, serves as primary key
    firmware_version DECIMAL(5, 2) NOT NULL,     -- Required firmware version
    network_type VARCHAR(50) NOT NULL,           -- Required network type
    created_at TIMESTAMP DEFAULT NOW(),          -- Timestamp when the record was created
    updated_at TIMESTAMP DEFAULT NOW(),   -- Updated automatically via trigger.
    timestamp BIGINT,
    flag SMALLINT DEFAULT 0,                     -- Additional column with default value, placed after timestamp

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

-- Attach a trigger to update the 'updated_at' field before updates.
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON parking.nbiot_device_settings
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Index for device lookups.
CREATE INDEX idx_nbiot_device_settings_device_id ON parking.nbiot_device_settings (device_id);
