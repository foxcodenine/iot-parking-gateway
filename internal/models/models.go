package models

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	up "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

// Global variables for the pgxpool connection and Upper ORM session.
var db *pgxpool.Pool
var dbSession up.Session

// Models struct will hold references to all database models, e.g., Device.
type Models struct {
	Device      Device
	RawDataLog  RawDataLog
	ActivityLog ActivityLog
	BeaconLog   BeaconLog
}

// New initializes the Models struct and sets up the Upper ORM session.
func New(conn *pgxpool.Pool) (Models, error) {

	// Assign the pgxpool connection pool to the global db variable.
	db = conn

	// Convert the pgxpool connection to a standard sql.DB.
	stdDB := stdlib.OpenDB(*db.Config().ConnConfig)

	// Initialize Upper ORM with the standard sql.DB connection.
	upperSession, err := postgresql.New(stdDB)
	if err != nil {
		return Models{}, err
	}

	// Assign the Upper ORM session to the global upper variable.
	dbSession = upperSession

	// Return an initialized Models struct with model references.
	return Models{}, nil
}
