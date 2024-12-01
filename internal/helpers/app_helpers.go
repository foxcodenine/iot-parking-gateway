package helpers

import (
	"context"
	"errors"

	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/middleware"
)

// GetUserFromContext retrieves the user claims from the context.
// It returns an error if the user claims cannot be found or are of the wrong type.
func GetUserFromContext(ctx context.Context) (*middleware.UserClaims, error) {
	userData, ok := ctx.Value(middleware.UserContextKey).(*middleware.UserClaims)
	if !ok || userData == nil {
		return nil, errors.New("user claims not found or wrong type in context")
	}
	return userData, nil
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
