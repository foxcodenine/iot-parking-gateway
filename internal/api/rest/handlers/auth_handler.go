package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
		Details:     fmt.Sprintf("User with ID %d and email '%s' logged in.", user.ID, user.Email),
	}

	// Push the audit log entry to the cache
	app.Cache.RPush("audit-logs", auditLogEntry)

	settings, _ := app.Cache.HGetAll("app:settings")

	// Respond with the token and user data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":           user.ID,
			"email":        user.Email,
			"access_level": user.AccessLevel,
			"enabled":      user.Enabled,
			"created_at":   user.CreatedAt,
		},
		"setting": settings,
	})
}
