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

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
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
	return "app.users"
}

var ErrDuplicateUser = errors.New("user with this email already exists")

// All retrieves all users either from the cache or the database if not cached.
func (u *User) GetAll() ([]*User, error) {
	var users []*User

	// Attempt to retrieve cached users
	cachedData, err := cache.AppCache.Get("db:users")
	if err != nil {
		helpers.LogError(err, "Failed to get users from cache")
		return nil, err // returning the error to handle it upstream
	}

	if cachedData != nil {
		// Asserting the type of cached data to []interface{}
		cachedUsers, ok := cachedData.([]interface{})
		if !ok {
			helpers.LogError(fmt.Errorf("cache data type mismatch: expected []interface{}, got %T", cachedData), "Cache data type mismatch")
			return nil, fmt.Errorf("cache data type mismatch: expected []interface{}, got %T", cachedData)
		}

		// Initialize slice to hold the converted user objects
		users = make([]*User, len(cachedUsers))
		for i, cachedUser := range cachedUsers {
			userMap, ok := cachedUser.(map[string]interface{})
			if !ok {
				helpers.LogError(fmt.Errorf("failed to assert type for user data: %T", cachedUser), "Error asserting type for cached user data")
				continue
			}

			user := &User{} // Create a new User instance
			// Map data from userMap to user struct fields
			if id, ok := userMap["id"].(float64); ok { // JSON numbers are by default float64
				user.ID = int(id)
			}
			if email, ok := userMap["email"].(string); ok {
				user.Email = email
			}
			if accessLevel, ok := userMap["access_level"].(float64); ok {
				user.AccessLevel = int(accessLevel)
			}
			if enabled, ok := userMap["enabled"].(bool); ok {
				user.Enabled = enabled
			}
			if createdAt, ok := userMap["created_at"].(string); ok { // Assuming stored as ISO8601 string
				user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
			}
			if updatedAt, ok := userMap["updated_at"].(string); ok {
				user.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
			}

			users[i] = user
		}
		// filter usrs
		return users, nil
	}

	// If not cached, fetch from the database
	collection := dbSession.Collection(u.TableName())
	err = collection.Find().All(&users)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users from database: %w", err)
	}

	// Cache the users after successful database fetch
	ttl, err := strconv.Atoi(os.Getenv("REDIS_DEFAULT_TTL"))
	if err != nil {
		helpers.LogError(err, "Failed to convert REDIS_DEFAULT_TTL to integer")
		ttl = 600 // Default TTL as a fallback
	}

	err = cache.AppCache.Set("db:users", users, ttl)
	if err != nil {
		helpers.LogError(err, "Failed to set users in cache")
	}

	return users, nil
}

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

	err = cache.AppCache.Delete("db:users")

	if err != nil {
		helpers.LogError(err, "failed to delete users from cache")
	}

	// Return the created user, including the ID
	return u, nil
}

func (u *User) Update(updatePassword bool) (*User, error) {
	// Validate required fields
	if strings.TrimSpace(u.Email) == "" {
		return nil, errors.New("email cannot be empty")
	}
	if !helpers.EmailRegex.MatchString(u.Email) {
		return nil, errors.New("invalid email format")
	}

	// Updating only the password if it's not empty, assumes other fields are managed separately
	if updatePassword {
		var err error
		u.Password, err = helpers.HashPassword(u.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
	}

	// Update timestamp for the user
	u.UpdatedAt = time.Now().UTC()

	// Execute the update operation
	collection := dbSession.Collection(u.TableName())
	err := collection.UpdateReturning(u)
	if err != nil {
		// Handle possible duplicate email error
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return nil, ErrDuplicateUser
		}
		// Wrap and return other errors
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Invalidate the cache after a successful update
	err = cache.AppCache.Delete("db:users")
	if err != nil {
		helpers.LogError(err, "failed to delete users from cache")
	}

	cache.AppCache.Set(fmt.Sprintf("app:user:logout:%d", u.ID), time.Now().Unix(), 86400)

	// Return the created user, including the ID
	return u, nil
}

func (u *User) Delete(userID int) error {
	// Reference the user collection
	collection := dbSession.Collection(u.TableName())

	// Check if the user exists before attempting to delete
	res := collection.Find(up.Cond{"id": userID})
	count, err := res.Count()
	if err != nil {
		// Log the error and return it
		helpers.LogError(err, fmt.Sprintf("Failed to count users with ID %d", userID))
		return fmt.Errorf("failed to verify user existence: %w", err)
	}

	// If no user is found, return an error
	if count == 0 {
		err := fmt.Errorf("user with ID %d does not exist", userID)
		helpers.LogError(err, "Delete operation failed: user not found")
		return err
	}

	// Attempt to delete the user
	if err := res.Delete(); err != nil {
		// Log the deletion error and return it
		helpers.LogError(err, fmt.Sprintf("Failed to delete user with ID %d", userID))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Invalidate the users cache after successful deletion
	err = cache.AppCache.Delete("db:users")
	if err != nil {
		// Log the cache deletion error but don't fail the operation
		helpers.LogError(err, "Failed to invalidate user cache after deletion")
	}

	cache.AppCache.Set(fmt.Sprintf("app:user:logout:%d", userID), time.Now().Unix(), 86400)

	return nil
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

// FindUserByID retrieves a user by their Email. Returns nil if the user is not found.
func (u *User) FindUserByID(userID int) (*User, error) {

	collection := dbSession.Collection(u.TableName())
	var user User
	err := collection.Find(up.Cond{"id": userID}).One(&user)
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

func (u *User) Upsert(newUser *User, updatePassword bool) (*User, error) {
	collection := dbSession.Collection(u.TableName())

	// Prepare the upsert query with positional placeholders
	sqlQuery := `
        INSERT INTO app.users (
            email, password, access_level, enabled, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6
        )
        ON CONFLICT (email) DO UPDATE SET
            access_level = EXCLUDED.access_level,
            enabled = EXCLUDED.enabled,
            updated_at = EXCLUDED.updated_at
    `

	// Prepare password logic based on whether we are updating the password
	hashedPassword := newUser.Password
	if updatePassword {
		var err error
		hashedPassword, err = helpers.HashPassword(newUser.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
	}

	// Prepare the parameter values
	params := []interface{}{
		newUser.Email,
		hashedPassword,
		newUser.AccessLevel,
		newUser.Enabled,
		time.Now().UTC(), // CreatedAt
		time.Now().UTC(), // UpdatedAt
	}

	// Execute the query
	_, err := collection.Session().SQL().Exec(sqlQuery, params...)
	if err != nil {
		return nil, helpers.WrapError(fmt.Errorf("failed to upsert user: %w", err))
	}

	// Invalidate the cache for users
	err = cache.AppCache.Delete("db:users")
	if err != nil {
		helpers.LogError(err, "Failed to invalidate users cache after upsert")
	}

	// Retrieve the upserted user to ensure return data is current
	updatedUser, err := u.FindUserByEmail(newUser.Email)
	if err != nil {
		return nil, helpers.WrapError(fmt.Errorf("failed to retrieve upserted user: %w", err))
	}

	return updatedUser, nil
}

func (u *User) GetRootUser() (*User, error) {

	collection := dbSession.Collection(u.TableName())
	var user User
	err := collection.Find(up.Cond{"access_level": 0}).One(&user)
	if err != nil {
		if err == up.ErrNoMoreRows {
			// No user found with the given email
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return &user, nil
}
