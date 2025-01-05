-- Devices table stores unique device identifiers and tracks their creation and update timestamps.
CREATE TABLE parking.devices (
    device_id VARCHAR(255) PRIMARY KEY,  -- IMEI, typically 15 digits.
    name VARCHAR(255) NULL,              -- Name of the device.
    network_type VARCHAR(50) NULL,       -- Type of network (e.g., NB-IoT, LoRa).
    firmware_version DECIMAL(5, 2) NULL, -- Firmware version of the device.
    latitude DECIMAL(9, 6) DEFAULT 0,    -- Latitude for geographic location.
    longitude DECIMAL(9, 6) DEFAULT 0,   -- Longitude for geographic location.
    beacons JSONB,                       -- JSONB data type for beacon information.
    happened_at TIMESTAMP NULL,          -- Timestamp of the device data capture.
    keepalive_at TIMESTAMP NULL,          
    settings_at TIMESTAMP NULL,          
    is_occupied BOOLEAN NULL,      -- Whether the space is occupied.
    is_allowed BOOLEAN DEFAULT FALSE,      -- Indicates if the device is allowed.
    is_blocked BOOLEAN DEFAULT FALSE,      -- Indicates if the device is blocked.
    is_hidden BOOLEAN DEFAULT FALSE,       -- Indicates if the device is hidden from view.
    created_at TIMESTAMP DEFAULT NOW(),  -- Set at record creation.
    updated_at TIMESTAMP DEFAULT NOW(),   -- Updated automatically via trigger.
    deleted_at TIMESTAMP DEFAULT NULL   -- Updated automatically via trigger.
);



-- Attach a trigger to update the 'updated_at' field before any update operation on 'devices'.
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON parking.devices
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Index for device lookups.
CREATE INDEX idx_devices_device_id ON parking.devices (device_id);
