CREATE TABLE parking.audit_logs (
    id SERIAL,
    user_id INTEGER NOT NULL,
    email VARCHAR(100) NOT NULL,
    access_level SMALLINT NOT NULL,
    
    -- Timestamps
    happened_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    action VARCHAR(50) NOT NULL,

    -- Optional fields related to the action
    entity VARCHAR(50),
    entity_id VARCHAR(255),
    url VARCHAR(100),
    ip_address VARCHAR(100),
    details TEXT,
    
    PRIMARY KEY (id, happened_at)
);

-- Convert 'audit_logs' to a TimescaleDB hypertable for time-based partitioning
SELECT create_hypertable('parking.audit_logs', 'happened_at');
