-- Enable TimescaleDB for time-series data management.
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- Create a function to automatically update the 'updated_at' timestamp.
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
