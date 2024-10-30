package models

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	db2 "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

var db *pgxpool.Pool
var upper db2.Session

// Models struct will hold references to all the database models.
type Models struct {
	Device Device
}

// New initializes the Models struct and sets up the upper/db session.
func New(conn *pgxpool.Pool) (Models, error) {
	// Assign the connection pool to the db variable.
	db = conn

	// Get a standard sql.DB from the pgxpool connection.
	stdDB := stdlib.OpenDB(*conn.Config().ConnConfig)

	// Initialize upper/db with the standard sql.DB.
	upperSession, err := postgresql.New(stdDB)
	if err != nil {
		return Models{}, err
	}

	upper = upperSession
	return Models{}, nil
}
