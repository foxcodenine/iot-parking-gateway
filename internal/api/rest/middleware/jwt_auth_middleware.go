package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/foxcodenine/iot-parking-gateway/internal/apptypes"
	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware creates a middleware for JWT authentication.
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the 'Authorization' header from the request
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Split the header into parts to separate the bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization header format must be 'Bearer {token}'", http.StatusUnauthorized)
			return
		}

		// Retrieve the JWT secret key from environment variables
		secret := os.Getenv("JWT_SECRET_KEY")
		tokenStr := parts[1] // The JWT token itself
		claims := &apptypes.UserClaims{}

		// Parse the token with claims
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			// This callback function returns the secret key for token verification
			return []byte(secret), nil
		})
		// Check for parsing errors or if the token is not valid
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Check if the user has been logged out
		logoutTimestampInterface, err := cache.AppCache.Get(fmt.Sprintf("app:user:logout:%d", claims.UserID))
		if err != nil {
			// Handle errors during retrieval from Redis
			http.Error(w, "Failed to retrieve logout timestamp from cache.", http.StatusInternalServerError)

			return
		}

		// Check if there's a logout timestamp for the user.
		if logoutTimestampInterface != nil {

			// Convert the retrieved value to int64. If conversion fails, report an error.
			logoutTimestampFloat64, ok := logoutTimestampInterface.(float64)
			if !ok {
				http.Error(w, "Invalid timestamp format", http.StatusInternalServerError)
				return
			}

			logoutTimestamp := int64(logoutTimestampFloat64)

			// Invalidate the token if it was issued before the logout timestamp.
			// This is triggered when an admin changes critical user account details like email or access level.
			// It ensures that users must re-authenticate to reflect these changes immediately.
			if claims.Timestamp < logoutTimestamp {
				http.Error(w, "Token is no longer valid", http.StatusUnauthorized)
				return
			}
		}

		// Add user data to context for access in subsequent handlers
		ctx := context.WithValue(r.Context(), apptypes.UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
