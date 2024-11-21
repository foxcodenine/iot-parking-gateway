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
