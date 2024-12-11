package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	up "github.com/upper/db/v4"
)

// Setting represents a configuration setting stored in the database.
type Setting struct {
	ID          int       `db:"id,omitempty" json:"id"`           // Unique identifier, automatically generated
	Key         string    `db:"key" json:"key"`                   // Unique key for the setting
	Val         string    `db:"val" json:"val"`                   // Value of the setting
	Description string    `db:"description" json:"description"`   // Description of what the setting controls or its purpose
	AccessLevel int       `db:"access_level" json:"access_level"` // Minimum access level required to view or edit this setting
	CreatedAt   time.Time `db:"created_at" json:"created_at"`     // Timestamp when the setting was created
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`     // Timestamp when the setting was last updated
	UpdatedBy   int       `db:"updated_by" json:"updated_by"`     // ID of the user who last updated this setting
}

// TableName returns the name of the table in the database, used by the ORM to map the struct to the database table.
func (Setting) TableName() string {
	return "app.settings"
}

// GetAll retrieves all settings from the database.
func (s *Setting) GetAll() ([]*Setting, error) {
	var settings []*Setting

	// If not cached, fetch from the database
	collection := dbSession.Collection(s.TableName())
	err := collection.Find().OrderBy("id").All(&settings)
	if err != nil {
		return nil, helpers.WrapError(fmt.Errorf("failed to retrieve settings from database: %w", err))
	}

	return settings, nil
}

// GetByKey retrieves a specific setting from the database based on its unique key.
func (s *Setting) GetByKey(key string) (*Setting, error) {
	collection := dbSession.Collection(s.TableName())

	var setting Setting

	// Fetch a single record where 'key' matches and 'deleted_at' is considered if soft deletes are implemented.
	err := collection.Find(up.Cond{"key": key}).One(&setting)
	if err != nil {
		if errors.Is(err, up.ErrNoMoreRows) {
			return nil, helpers.WrapError(errors.New("setting not found"))
		}
		return nil, helpers.WrapError(fmt.Errorf("failed to retrieve setting: %w", err))
	}

	return &setting, nil
}

// Create inserts a new setting into the database and returns the created setting.
func (s *Setting) Create(newSetting *Setting) (*Setting, error) {
	collection := dbSession.Collection(s.TableName())

	// Set current time for CreatedAt and UpdatedAt
	now := time.Now().UTC()
	newSetting.CreatedAt = now
	newSetting.UpdatedAt = now

	_, err := collection.Insert(newSetting)
	if err != nil {
		// Check if the error is a duplicate key violation (PostgreSQL SQL state 23505)
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			var ErrDuplicateKey = helpers.WrapError(errors.New("a setting with this key already exists"))
			return nil, ErrDuplicateKey
		}

		// Return any other errors with additional context
		return nil, helpers.WrapError(fmt.Errorf("failed to create setting: %w", err))
	}

	// Attempt to cache the new setting in Redis
	err = cache.AppCache.HSet("app:settings", newSetting.Key, newSetting.Val)
	if err != nil {
		// Log the error without disrupting the main flow as caching failure is often not critical
		helpers.LogError(err, "Failed to cache new setting in Redis")
	}

	return newSetting, nil
}

// Upsert creates a new setting or updates it if it already exists based on its unique key.
func (s *Setting) Upsert(newSetting *Setting) (*Setting, error) {
	collection := dbSession.Collection(s.TableName())

	// Prepare the upsert query with positional placeholders
	sqlQuery := `
        INSERT INTO app.settings (
            key, val, description, access_level, created_at, updated_at, updated_by
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7
        )
        ON CONFLICT (key) DO UPDATE SET
            val = EXCLUDED.val,
            description = EXCLUDED.description,
            access_level = EXCLUDED.access_level,
            updated_at = EXCLUDED.updated_at,
            updated_by = EXCLUDED.updated_by;
    `

	// Prepare the parameter values
	params := []interface{}{
		newSetting.Key,
		newSetting.Val,
		newSetting.Description,
		newSetting.AccessLevel,
		time.Now().UTC(),
		time.Now().UTC(),
		newSetting.UpdatedBy,
	}

	// Execute the query
	_, err := collection.Session().SQL().Exec(sqlQuery, params...)
	if err != nil {
		return nil, helpers.WrapError(fmt.Errorf("failed to upsert setting: %w", err))
	}

	// Retrieve the upserted setting to ensure return data is current
	updatedSetting, err := s.GetByKey(newSetting.Key)
	if err != nil {
		return nil, helpers.WrapError(fmt.Errorf("failed to retrieve upserted setting: %w", err))
	}

	// Attempt to update the setting in Redis
	err = cache.AppCache.HSet("app:settings", newSetting.Key, updatedSetting.Val)
	if err != nil {
		// Log the error without disrupting the main flow as caching failure is often not critical
		helpers.LogError(err, "Failed to update setting in Redis")
	}

	return updatedSetting, nil
}

// UpdateByKey updates specific fields of a setting identified by its unique key.
func (s *Setting) UpdateByKey(key string, updatedFields map[string]interface{}) (*Setting, error) {
	collection := dbSession.Collection(s.TableName())

	// Set the time of update to ensure it's always current
	updatedFields["updated_at"] = time.Now().UTC()

	// Check if the setting exists by counting matching rows
	res := collection.Find(up.Cond{"key": key})
	count, err := res.Count()
	if err != nil {
		return nil, helpers.WrapError(fmt.Errorf("error checking setting existence: %w", err))
	}

	// If no matching setting is found, return an error indicating so
	if count == 0 {
		return nil, helpers.WrapError(fmt.Errorf("setting with key '%s' not found", key))
	}

	// Perform the update on the found record
	err = res.Update(updatedFields)
	if err != nil {
		return nil, helpers.WrapError(fmt.Errorf("error updating setting: %w", err))
	}

	// Retrieve the updated setting to ensure return data is current
	updatedSetting, err := s.GetByKey(key)
	if err != nil {
		return nil, helpers.WrapError(err)
	}

	// Attempt to update the setting in Redis
	err = cache.AppCache.HSet("app:settings", key, updatedSetting.Val)
	if err != nil {
		// Log the error without disrupting the main flow as caching failure is often not critical
		helpers.LogError(err, "Failed to update setting in Redis")
	}

	return updatedSetting, nil
}
