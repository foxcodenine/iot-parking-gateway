package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
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

		helpers.LogError(err, "Error finding user:")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Invalid email or password!", http.StatusUnauthorized)
		return
	}
	// Verify the password
	if !helpers.CheckPasswordHash(payload.Password, user.Password) {
		http.Error(w, "Invalid email or password!!", http.StatusUnauthorized)
		return
	}

	// Generate a token (e.g., JWT or session token)
	token, err := user.GenerateToken() // Assuming `GenerateToken` is implemented
	if err != nil {
		helpers.LogError(err, "Error generating token:")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

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
	})
}
