CREATE TABLE app.settings (
    id SERIAL PRIMARY KEY,
    key VARCHAR(255) UNIQUE NOT NULL,
    val VARCHAR(255) NOT NULL,
    description TEXT,
    access_level INT NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    updated_by INT    
);

-- Attach a trigger to update the 'updated_at' field before updates.
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON app.settings
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
