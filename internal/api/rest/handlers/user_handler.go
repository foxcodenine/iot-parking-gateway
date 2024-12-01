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

func (u *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	users, err := app.Models.User.All()

	if err != nil {
		// Log the error and send an HTTP 500 Internal Server Error response
		http.Error(w, "Unable to retrieve users.", http.StatusInternalServerError)
		helpers.LogError(err, "Failed to retrieve users from the database.")
		return
	}

	userData, err := app.GetUserFromContext(r.Context())
	if err != nil {
		app.ErrorLog.Printf("Authentication error: %v", err)
		http.Error(w, "Authentication error.", http.StatusUnauthorized)
		return
	}

	// Initialize the slice for filtered users
	filteredUsers := make([]*models.User, 0)

	// Apply filtering based on user access level
	if userData.AccessLevel == 0 {
		// Root user, access level 0, can see all accounts
		filteredUsers = users
	} else {
		// Non-root users, exclude users with root access
		for _, user := range users {
			if user.AccessLevel > 0 { // Exclude root level users
				filteredUsers = append(filteredUsers, user)
			}
		}
	}

	// Response structure with a success message and user data
	response := map[string]interface{}{
		"message": "Users retrieved successfully.",
		"users":   filteredUsers,
	}

	// Set content type to application/json before writing the status or data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // HTTP 200 OK for a successful GET request

	// Encode the response as JSON and handle any encoding errors
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode users data as JSON.", http.StatusInternalServerError)
		helpers.LogError(err, "Error encoding users data as JSON:")
	}

}

func (u *UserHandler) Store(w http.ResponseWriter, r *http.Request) {

	userData, err := app.GetUserFromContext(r.Context())
	if err != nil {
		app.ErrorLog.Printf("Authentication error: %v", err)
		http.Error(w, "Authentication error.", http.StatusUnauthorized)
		return
	}

	// Check if the user has permission to create a new user
	if userData.AccessLevel > 1 {
		http.Error(w, "You do not have the necessary permissions to perform this action.", http.StatusForbidden)
		return
	}

	// Parse and validate input from the API
	type Request struct {
		Email       string `json:"email"`
		Password1   string `json:"password1"`
		Password2   string `json:"password2"`
		AccessLevel int    `json:"access_level"`
	}

	var req Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondWithError(w, err, "Failed to create user.", http.StatusInternalServerError)
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.Email) == "" {
		http.Error(w, "Email cannot be empty.", http.StatusBadRequest)
		return
	}
	if !helpers.EmailRegex.MatchString(req.Email) { // Assuming emailRegex is defined globally for email validation
		http.Error(w, "Invalid email format.", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Password1) == "" || strings.TrimSpace(req.Password2) == "" {
		http.Error(w, "Passwords cannot be empty.", http.StatusBadRequest)
		return
	}
	if len(req.Password1) < 6 {
		http.Error(w, "Password must be at least 6 characters long.", http.StatusBadRequest)
		return
	}

	// Validate that passwords match
	if req.Password1 != req.Password2 {
		http.Error(w, "Passwords do not match.", http.StatusBadRequest)
		return
	}

	// Validate access level
	if req.AccessLevel < 0 || req.AccessLevel > 3 { // Example range validation
		http.Error(w, "Invalid access level.", http.StatusBadRequest)
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
			http.Error(w, "User with this email already exists.", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create user.", http.StatusInternalServerError)
		helpers.LogError(err, "Error creating user:")
		return
	}

	// Respond with a success message and the created user's data (excluding sensitive info)
	response := map[string]interface{}{
		"message": "User created successfully.",
		"user":    createdUser, // Directly using the createdUser struct
	}

	// Respond with the created user's data (excluding sensitive info)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle error if the JSON encoding fails
		http.Error(w, "Failed to send response.", http.StatusInternalServerError)
		helpers.LogError(err, "Error encoding response:")
	}

}
