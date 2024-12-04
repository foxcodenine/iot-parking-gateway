package apptypes

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims defines the structure of the JWT claims used in the application.
type UserClaims struct {
	jwt.RegisteredClaims        // Embedding the standard jwt claims
	AccessLevel          int    `json:"access_level"` // User access level
	Email                string `json:"email"`        // User's email address
	Timestamp            int64  `json:"timestamp"`    // The timestamp when the JWT was issued
	UserID               int    `json:"user_id"`      // Unique identifier for the user
}

// Define a custom type to ensure that context keys are unique across the application
type contextKey int

// Declare constants for the keys using the custom type
const (
	UserContextKey contextKey = iota // iota increments automatically, userContextKey will be 0
	// Add other keys here, each will have a unique integer value
)

// GetUserFromContext retrieves the user claims from the context.
// It returns an error if the user claims cannot be found or are of the wrong type.
func GetUserFromContext(ctx context.Context) (*UserClaims, error) {
	userData, ok := ctx.Value(UserContextKey).(*UserClaims)
	if !ok || userData == nil {
		return nil, errors.New("user claims not found or wrong type in context")
	}
	return userData, nil
}
