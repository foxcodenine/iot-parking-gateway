CREATE TABLE parking.users (
    id SERIAL PRIMARY KEY,                  -- Unique identifier for each user
    username VARCHAR(255) NOT NULL UNIQUE,  -- Username, must be unique
    password VARCHAR(255) NOT NULL,         -- Password for the user
    access_level INTEGER NOT NULL DEFAULT 3,-- Access level (e.g., roles or permissions)
    enabled BOOLEAN NOT NULL DEFAULT TRUE,   -- Status to indicate if the user is enabled
    created_at TIMESTAMP DEFAULT NOW(),  -- Set at record creation.
    updated_at TIMESTAMP DEFAULT NOW()   -- Updated automatically via trigger.
);

-- Attach a trigger to update the 'updated_at' field before updates.
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON parking.users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();