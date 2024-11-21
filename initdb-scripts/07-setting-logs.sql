-- NB-IoT setting logs table tracks configuration changes.
CREATE TABLE parking.nbiot_setting_logs (
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
    nb_iot_imsi VARCHAR(15),

    PRIMARY KEY (id, happened_at)
);
-- Convert 'nbiot_setting_logs' to a TimescaleDB hypertable for optimized performance.
SELECT create_hypertable('parking.nbiot_setting_logs', 'happened_at');

-- Indexes for improved query performance.
CREATE INDEX idx_nbiot_setting_logs_device_id ON parking.nbiot_setting_logs (device_id);
CREATE INDEX idx_nbiot_setting_logs_happened_at ON parking.nbiot_setting_logs (happened_at);
