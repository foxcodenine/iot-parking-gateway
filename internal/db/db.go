package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// OpenDB initializes a connection pool to the PostgreSQL database and returns a pgxpool.Pool instance.
func OpenDB() (*pgxpool.Pool, error) {

	// Format the DSN (Data Source Name) using environment variables.
	var dbPort string

	switch os.Getenv("GO_ENV") {
	case "production":
		dbPort = os.Getenv("DB_PORT")
	default:
		dbPort = os.Getenv("DB_PORT_EX")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		dbPort,
		os.Getenv("DB_NAME"),
	)

	// Parse the DSN into a configuration for the pgx connection pool.

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// Create a new connection pool using the parsed configuration.
	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure the connection is established.
	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	// Return the connection pool if all steps succeeded.
	return db, nil
}
