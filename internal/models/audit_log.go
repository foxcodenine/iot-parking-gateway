package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

// AuditLog represents a single entry in the audit_logs table.
type AuditLog struct {
	ID          int    `db:"id,omitempty" json:"id,omitempty"` // Unique identifier
	UserID      int    `db:"user_id" json:"user_id"`           // ID of the user who performed the action
	Email       string `db:"email" json:"email"`               // Email of the user who performed the action
	AccessLevel int    `db:"access_level" json:"access_level"` // Access level of the user

	// Timestamps
	HappenedAt time.Time `db:"happened_at" json:"happened_at"` // Timestamp of when the action occurred

	Action string `db:"action" json:"action"` // Action performed (e.g., login, delete)

	// Optional Contextual Information (no longer pointers)
	Entity    string `db:"entity" json:"entity,omitempty"`         // Name of the entity affected
	EntityID  string `db:"entity_id" json:"entity_id,omitempty"`   // ID of the entity affected
	URL       string `db:"url" json:"url,omitempty"`               // URL accessed
	IPAddress string `db:"ip_address" json:"ip_address,omitempty"` // IP address of the user
	Details   string `db:"details" json:"details,omitempty"`       // Additional details about the action
}

// TableName returns the full table name for the AuditLog model in PostgreSQL.
func (a *AuditLog) TableName() string {
	return "app.audit_logs"
}

// NewAuditLog constructs an AuditLog object from a provided map of data.
// It handles data type conversions and populates the fields accordingly.
func NewAuditLog(data map[string]interface{}) (*AuditLog, error) {
	// Parse the required fields
	userID, ok := data["user_id"].(float64)
	if !ok {
		return nil, helpers.WrapError(errors.New("invalid or missing 'user_id' field"))
	}

	happenedAtStr, ok := data["happened_at"].(string)
	if !ok {
		return nil, helpers.WrapError(errors.New("invalid or missing 'happened_at' field"))
	}
	happenedAt, err := time.Parse(time.RFC3339, happenedAtStr)
	if err != nil {
		return nil, helpers.WrapError(fmt.Errorf("invalid 'happened_at' format: %v", err))
	}

	email, ok := data["email"].(string)
	if !ok || email == "" {
		return nil, helpers.WrapError(errors.New("invalid or missing 'email' field"))
	}

	accessLevel, ok := data["access_level"].(float64)
	if !ok {
		return nil, helpers.WrapError(errors.New("invalid or missing 'access_level' field"))
	}

	action, ok := data["action"].(string)
	if !ok || action == "" {
		return nil, helpers.WrapError(errors.New("invalid or missing 'action' field"))
	}

	// Construct the AuditLog object
	auditLog := &AuditLog{
		UserID:      int(userID),
		HappenedAt:  happenedAt,
		Email:       email,
		AccessLevel: int(accessLevel),
		Action:      action,
	}

	// Parse optional fields
	if entity, ok := data["entity"].(string); ok {
		auditLog.Entity = entity
	}

	if entityID, ok := data["entity_id"].(string); ok {

		auditLog.EntityID = entityID
	}

	if ipAddress, ok := data["ip_address"].(string); ok {
		auditLog.IPAddress = ipAddress
	}

	if url, ok := data["url"].(string); ok {
		auditLog.URL = url
	}

	if details, ok := data["details"].(string); ok {
		auditLog.Details = details
	}

	return auditLog, nil
}

func (a *AuditLog) BulkInsert(auditLog []AuditLog) error {
	// Exit early if there are no records to insert
	if len(auditLog) == 0 {
		return nil
	}

	// Determine the number of fields to be inserted for each log
	numFields := 10 // Adjust this based on the actual number of columns you have in your table

	// Prepare slices for SQL values and arguments.
	values := make([]string, 0, len(auditLog))
	args := make([]interface{}, 0, len(auditLog)*numFields) // Adjust the argument count based on the number of columns

	for i, log := range auditLog {
		// Create a placeholder for each record with indexed arguments
		placeholders := make([]string, numFields)
		for j := range placeholders {
			placeholders[j] = fmt.Sprintf("$%d", i*numFields+j+1)
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))

		// Append the actual values for each placeholder in the same order as the columns
		args = append(args,
			log.UserID, log.Email, log.AccessLevel, log.HappenedAt, log.Action, log.Entity, log.EntityID,
			log.URL, log.IPAddress, log.Details,
		)
	}

	// Construct the SQL statement by joining the placeholders for each record
	query := fmt.Sprintf("INSERT INTO %s (user_id, email, access_level, happened_at, action, entity, entity_id, url, ip_address, details) VALUES %s",
		a.TableName(), strings.Join(values, ", "))

	// Execute the constructed query with the arguments
	_, err := dbSession.SQL().Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk insert for setting logs: %w", err)
	}

	return nil
}
