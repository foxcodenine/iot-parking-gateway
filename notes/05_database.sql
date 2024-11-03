CREATE SCHEMA parking;

CREATE TABLE parking.devices (
    device_id VARCHAR(15) PRIMARY KEY,      -- IMEI is typically 15 digits
    created_at TIMESTAMP DEFAULT NOW(),     -- Automatically set when the record is created
    updated_at TIMESTAMP DEFAULT NOW()      -- Will be managed by a trigger for automatic updates
);


CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();  -- Update the updated_at column to the current timestamp
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON parking.devices
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- ---------------------------------------------------------------------

CREATE TABLE parking.raw_data_logs (
    uudi UUID DEFAULT gen_random_uuid(),
    hex_data TEXT NOT NULL,
    happened_at TIMESTAMP NOT NULL,         -- UTC, no timezone
    created_at TIMESTAMP DEFAULT NOW(),     -- UTC timestamp of saving
    PRIMARY KEY (uudi, happened_at)
);

-- Convert to a TimescaleDB hypertable
SELECT create_hypertable('parking.raw_data_logs', 'happened_at');

-- ---------------------------------------------------------------------

