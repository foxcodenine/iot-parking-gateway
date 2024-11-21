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

-- Index for device lookups.
CREATE INDEX idx_devices_device_id ON parking.devices (device_id);
