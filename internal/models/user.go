package models

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	up "github.com/upper/db/v4"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

// User represents a user in the database with login credentials and permissions.
type User struct {
	ID          int       `db:"id,omitempty" json:"id"`           // Unique identifier
	Email       string    `db:"email" json:"email"`               // Email for login
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

var ErrDuplicateUser = errors.New("user with this email already exists")

func (u *User) Create() (*User, error) {
	// Validate required fields
	if strings.TrimSpace(u.Email) == "" {
		return nil, errors.New("email cannot be empty")
	}
	if !helpers.EmailRegex.MatchString(u.Email) {
		return nil, errors.New("invalid email format")
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

	// Insert user into the database and get the generated ID
	collection := dbSession.Collection(u.TableName())
	err = collection.InsertReturning(u)
	if err != nil {
		// Check for duplicate email (SQLSTATE 23505 = unique violation)
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return nil, ErrDuplicateUser
		}

		// Wrap and return other errors
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Return the created user, including the ID
	return u, nil
}

// FindUserByEmail retrieves a user by their Email. Returns nil if the user is not found.
func (u *User) FindUserByEmail(email string) (*User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email cannot be empty")
	}

	collection := dbSession.Collection(u.TableName())
	var user User
	err := collection.Find(up.Cond{"email": email}).One(&user)
	if err != nil {
		if err == up.ErrNoMoreRows {
			// No user found with the given email
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return &user, nil
}

func (u *User) GenerateToken() (string, error) {
	// Load expiration time from environment
	ttlStr := os.Getenv("JWT_EXPIRATION_TIME")
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		return "", fmt.Errorf("invalid JWT_EXPIRATION_TIME: %v", err)
	}

	// Define token claims
	claims := jwt.MapClaims{
		"user_id":      u.ID,
		"email":        u.Email,
		"access_level": u.AccessLevel,
		"timestamp":    time.Now().Unix(),
		"exp":          time.Now().Add(time.Second * time.Duration(ttl)).Unix(), // 24-hour expiration
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
