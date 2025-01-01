-- NB-IoT device settings table stores device configurations.
CREATE TABLE parking.sigfox_device_settings (
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
    radar_car_delta_th INTEGER,
    downlink_en_7_bits_repeated_occupancy_period_mins SMALLINT
);

-- Attach a trigger to update the 'updated_at' field before updates.
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON parking.sigfox_device_settings
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Index for device lookups.
CREATE INDEX idx_sigfox_device_settings_device_id ON parking.sigfox_device_settings (device_id);
