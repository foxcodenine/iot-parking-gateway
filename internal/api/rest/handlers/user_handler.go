package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

type UserHandler struct {
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	// Parse and validate input from the API
	type Request struct {
		Email       string `json:"email"`
		Password1   string `json:"password1"`
		Password2   string `json:"password2"`
		AccessLevel int    `json:"access_level"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondWithError(w, err, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.Email) == "" {
		http.Error(w, "Email cannot be empty", http.StatusBadRequest)
		return
	}
	if !helpers.EmailRegex.MatchString(req.Email) { // Assuming emailRegex is defined globally for email validation
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Password1) == "" || strings.TrimSpace(req.Password2) == "" {
		http.Error(w, "Passwords cannot be empty", http.StatusBadRequest)
		return
	}
	if len(req.Password1) < 6 {
		http.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	// Validate that passwords match
	if req.Password1 != req.Password2 {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// Validate access level
	if req.AccessLevel < 0 || req.AccessLevel > 3 { // Example range validation
		http.Error(w, "Invalid access level", http.StatusBadRequest)
		return
	}

	// Create a new user instance
	newUser := &models.User{
		Email:       req.Email,
		Password:    req.Password1, // Password will be hashed in the `Create` method
		AccessLevel: req.AccessLevel,
		Enabled:     true,
	}

	// Attempt to create the user
	createdUser, err := newUser.Create()
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUser) {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		helpers.LogError(err, "Error creating user:")
		return
	}

	// Respond with the created user's data (excluding sensitive info)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":           createdUser.ID,
		"email":        createdUser.Email,
		"access_level": createdUser.AccessLevel,
		"enabled":      createdUser.Enabled,
		"created_at":   createdUser.CreatedAt,
	})

}
