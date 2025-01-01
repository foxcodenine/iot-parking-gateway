-- NB-IoT setting logs table tracks configuration changes.
CREATE TABLE parking.lora_setting_logs (
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
    radar_car_uncal_lo_th INTEGER,
    radar_car_uncal_hi_th INTEGER,
    radar_car_delta_th INTEGER,
    mag_car_lo INTEGER,
    mag_car_hi INTEGER,
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

    lora_data_rate SMALLINT,
    lora_retries SMALLINT,

    PRIMARY KEY (id, happened_at)
);
-- Convert 'lora_setting_logs' to a TimescaleDB hypertable for optimized performance.
SELECT create_hypertable('parking.lora_setting_logs', 'happened_at');

-- Indexes for improved query performance.
CREATE INDEX idx_lora_setting_logs_device_id ON parking.lora_setting_logs (device_id);
CREATE INDEX idx_lora_setting_logs_happened_at ON parking.lora_setting_logs (happened_at);
