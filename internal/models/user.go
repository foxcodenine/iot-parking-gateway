package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

// User represents a user in the database with login credentials and permissions.
type User struct {
	ID          int       `db:"id,omitempty" json:"id"`           // Unique identifier
	Username    string    `db:"username" json:"username"`         // Username for login
	Password    string    `db:"password" json:"-"`                // Hashed password (not serialized in JSON)
	AccessLevel int       `db:"access_level" json:"access_level"` // User's access level
	Enabled     bool      `db:"enabled" json:"enabled"`           // Whether the user is active
	CreatedAt   time.Time `db:"created_at" json:"created_at"`     // Timestamp for when the user was created
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`     // Timestamp for the last update
}

// TableName returns the full table name for the User model in PostgreSQL.
func (u *User) TableName() string {
	return "parking.users"
}

var ErrDuplicateUser = errors.New("user with this username already exists")

func (u *User) Create() (*User, error) {
	// Validate required fields
	if strings.TrimSpace(u.Username) == "" {
		return nil, errors.New("username cannot be empty")
	}
	if strings.TrimSpace(u.Password) == "" {
		return nil, errors.New("password cannot be empty")
	}

	// Set timestamps
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now

	// Hash the password
	var err error
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	fmt.Println(u)

	// Insert user into the database and get the generated ID
	collection := dbSession.Collection(u.TableName())
	err = collection.InsertReturning(u)
	if err != nil {
		// Check for duplicate username (SQLSTATE 23505 = unique violation)
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return nil, ErrDuplicateUser
		}

		// Wrap and return other errors
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Return the created user, including the ID
	return u, nil
}
