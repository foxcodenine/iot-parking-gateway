package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

type AuthHandler struct {
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the JSON payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the email and password fields
	payload.Email = strings.TrimSpace(payload.Email)
	payload.Password = strings.TrimSpace(payload.Password)

	if payload.Email == "" {
		http.Error(w, "Email cannot be empty", http.StatusBadRequest)
		return
	}

	if !helpers.EmailRegex.MatchString(payload.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if payload.Password == "" {
		http.Error(w, "Password cannot be empty", http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.FindUserByEmail(payload.Email)
	if err != nil {

		helpers.RespondWithError(w, err, "Internal server error. Error finding user.", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Invalid email or password!", http.StatusUnauthorized)
		return
	}

	// Verify the password
	if !helpers.CheckPasswordHash(payload.Password, user.Password) {
		http.Error(w, "Invalid email or password!", http.StatusUnauthorized)
		return
	}

	if !user.Enabled {
		http.Error(w, "User account is disabled. Please contact support.", http.StatusForbidden)
		return
	}

	// Generate a token (e.g., JWT or session token)
	token, err := user.GenerateToken() // Assuming `GenerateToken` is implemented
	if err != nil {
		helpers.RespondWithError(w, err, "Internal server error. Error generating token.", http.StatusInternalServerError)
		return
	}

	// Create the audit log entry
	auditLogEntry := models.AuditLog{
		UserID:      user.ID,
		Email:       user.Email,
		AccessLevel: user.AccessLevel,
		HappenedAt:  time.Now().UTC(),
		Action:      "LOGIN",
		URL:         r.URL.Path,
		IPAddress:   getClientIP(r),
		Details:     fmt.Sprintf("User with ID %d and email '%s' successfully logged.", user.ID, user.Email),
	}

	// Push the audit log entry to the cache
	app.Cache.RPush("logs:audit-logs", auditLogEntry)

	settings, err := app.Cache.HGetAll("app:settings")

	if err != nil {
		helpers.RespondWithError(w, err, "Failed to retrieve application settings.", http.StatusInternalServerError)
		return
	}

	favoritesData, err := cache.AppCache.HGet("app:user:favorites", fmt.Sprintf("%d", user.ID))
	favoritesDeviceIDs := []string{}

	if err != nil {
		helpers.LogError(err, "Error fetching favorites from Redis")
	} else if favoritesData != nil {

		// Convert []interface{} to []string
		if data, ok := favoritesData.([]interface{}); ok {
			for _, v := range data {
				if str, ok := v.(string); ok {
					favoritesDeviceIDs = append(favoritesDeviceIDs, str)
				} else {
					helpers.LogError(fmt.Errorf("unexpected type in favorites data: %T", v), "Error converting favorite to string")
				}
			}
		} else {
			helpers.LogError(fmt.Errorf("unexpected type for favoritesData: %T", favoritesData), "Error converting favorites from Redis")
		}
	}

	// Respond with the token and user data
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":           user.ID,
			"email":        user.Email,
			"access_level": user.AccessLevel,
			"enabled":      user.Enabled,
			"created_at":   user.CreatedAt,
			"favorites":    favoritesDeviceIDs,
		},
		"settings": settings,
	})

	if err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response.", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userData, err := app.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Authentication error.", http.StatusUnauthorized)
		return
	}

	// Create the audit log entry
	auditLogEntry := models.AuditLog{
		UserID:      userData.UserID,
		Email:       userData.Email,
		AccessLevel: userData.AccessLevel,
		HappenedAt:  time.Now().UTC(),
		Action:      "LOGOUT",
		URL:         r.URL.Path,
		IPAddress:   getClientIP(r),
		Details:     fmt.Sprintf("User with ID %d and email '%s' successfully logged out.", userData.UserID, userData.Email),
	}

	// Push the audit log entry to the cache
	app.Cache.RPush("logs:audit-logs", auditLogEntry)
}
