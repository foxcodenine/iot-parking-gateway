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

   deep_sleep_time_1 INTEGER,
    action_before_1 SMALLINT,
    action_after_1 SMALLINT,    
    deep_sleep_time_2 INTEGER,
    action_before_2 SMALLINT,
    action_after_2 SMALLINT,
    deep_sleep_time_3 INTEGER,
    action_before_3 SMALLINT,
    action_after_3 SMALLINT,
    deep_sleep_time_4 INTEGER,
    action_before_4 SMALLINT,
    action_after_4 SMALLINT,
    deep_sleep_time_5 INTEGER,
    action_before_5 SMALLINT,
    action_after_5 SMALLINT,
    deep_sleep_time_6 INTEGER,
    action_before_6 SMALLINT,
    action_after_6 SMALLINT,
    deep_sleep_time_7 INTEGER,
    action_before_7 SMALLINT,
    action_after_7 SMALLINT,
    deep_sleep_time_8 INTEGER,
    action_before_8 SMALLINT,
    action_after_8 SMALLINT,
    deep_sleep_time_9 INTEGER,
    action_before_9 SMALLINT,
    action_after_9 SMALLINT,
    deep_sleep_time_10 INTEGER,
    action_before_10 SMALLINT,
    action_after_10 SMALLINT,

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
