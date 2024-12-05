package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
}

func (u *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	users, err := app.Models.User.GetAll()

	if err != nil {
		// Log the error and send an HTTP 500 Internal Server Error response
		helpers.RespondWithError(w, err, "Unable to retrieve users.", http.StatusInternalServerError)
		return
	}

	userData, err := app.GetUserFromContext(r.Context())
	if err != nil {
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
		helpers.RespondWithError(w, err, "Failed to encode users data as JSON.", http.StatusInternalServerError)
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
	if req.AccessLevel < 1 || req.AccessLevel > 3 { // Example range validation
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
		helpers.RespondWithError(w, err, "Failed to create user.", http.StatusInternalServerError)
		return
	}

	// Respond with a success message and the created user's data (excluding sensitive info)
	response := map[string]interface{}{
		"message": "User created successfully.",
		"user":    createdUser, // Directly using the createdUser struct
	}

	app.PushAuditToCache(*userData, "CREATE", "user", fmt.Sprintf("%d", newUser.ID), r, fmt.Sprintf("Created user with ID %d and email '%s'.", newUser.ID, newUser.Email))

	// Respond with the created user's data (excluding sensitive info)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Handle error if the JSON encoding fails
		helpers.RespondWithError(w, err, "Failed to encoding and send response.", http.StatusInternalServerError)
	}
}

func (u *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userData, err := app.GetUserFromContext(r.Context())
	if err != nil {
		app.ErrorLog.Printf("Authentication error: %v", err)
		http.Error(w, "Authentication error.", http.StatusUnauthorized)
		return
	}

	// Check if the user has permission to update a user
	if userData.AccessLevel > 1 {
		http.Error(w, "You do not have the necessary permissions to perform this action.", http.StatusForbidden)
		return
	}

	// Fetch user ID from URL parameters
	userIDStr := chi.URLParam(r, "id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID.", http.StatusBadRequest)
		return
	}

	// Parse and validate input from the API
	type Request struct {
		Email         string `json:"email"`
		Password1     string `json:"password1"`
		Password2     string `json:"password2"`
		AccessLevel   int    `json:"access_level"`
		AdminPassword string `json:"admin_password"`
		Enabled       bool   `json:"enabled"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}

	// Verify admin password
	adminUser, err := app.Models.User.FindUserByID(userData.UserID)
	if err != nil || !helpers.CheckPasswordHash(req.AdminPassword, adminUser.Password) {
		http.Error(w, "Authentication failed: incorrect admin credentials.", http.StatusForbidden)
		return
	}

	// Find the user to update
	user, err := app.Models.User.FindUserByID(userID)
	if user == nil || err != nil {
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}

	// Apply updates from the request to the user model
	if strings.TrimSpace(req.Email) != "" {
		if !helpers.EmailRegex.MatchString(req.Email) {
			http.Error(w, "Invalid email format.", http.StatusBadRequest)
			return
		}
		user.Email = strings.TrimSpace(req.Email)
	}

	updatePassword := false

	if strings.TrimSpace(req.Password1) != "" {
		if req.Password1 != req.Password2 {
			http.Error(w, "Passwords do not match.", http.StatusBadRequest)
			return
		}
		if len(req.Password1) < 6 {
			http.Error(w, "Password must be at least 6 characters long.", http.StatusBadRequest)
			return
		}
		updatePassword = true
		user.Password = req.Password1 // Password will be hashed in the Update method
	}

	if req.AccessLevel >= 0 && req.AccessLevel <= 3 { // Assuming 0-3 are valid access levels
		if req.AccessLevel > 0 {
			user.AccessLevel = req.AccessLevel
		}
	} else {
		http.Error(w, "Invalid access level.", http.StatusBadRequest)
		return
	}

	user.Enabled = req.Enabled

	// Update the user in the database
	updatedUser, err := user.Update(updatePassword)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUser) {
			http.Error(w, "User with this email already exists.", http.StatusConflict)
			return
		}
		helpers.RespondWithError(w, err, "Failed to update user.", http.StatusInternalServerError)
		return
	}

	app.PushAuditToCache(*userData, "UPDATE", "user", fmt.Sprintf("%d", userID), r, fmt.Sprintf("Updated user with ID %d and email '%s'.", userID, updatedUser.Email))

	// Respond with success
	response := map[string]interface{}{
		"message": "User updated successfully.",
		"user":    updatedUser,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response.", http.StatusInternalServerError)
	}
}

func (u *UserHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	userData, err := app.GetUserFromContext(r.Context())
	if err != nil {
		app.ErrorLog.Printf("Authentication error: %v", err)
		http.Error(w, "Authentication error.", http.StatusUnauthorized)
		return
	}

	// Check if the user has permission to delete a user
	if userData.AccessLevel > 1 {
		http.Error(w, "You do not have the necessary permissions to perform this action.", http.StatusForbidden)
		return
	}

	// Fetch user ID from URL parameters
	userIDStr := chi.URLParam(r, "id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID.", http.StatusBadRequest)
		return
	}

	// Parse and validate input from the API
	type Request struct {
		AdminPassword string `json:"admin_password"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}

	// Verify admin password
	adminUser, err := app.Models.User.FindUserByID(userData.UserID)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to retrieve admin user for verification.", http.StatusInternalServerError)
		return
	}

	if !helpers.CheckPasswordHash(req.AdminPassword, adminUser.Password) {
		http.Error(w, "Authentication failed: incorrect admin credentials.", http.StatusForbidden)
		return
	}

	// Retrieve the user to delete
	user, err := app.Models.User.FindUserByID(userID)
	if err != nil {
		helpers.RespondWithError(w, err, fmt.Sprintf("Failed to retrieve user with ID %d.", userID), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, fmt.Sprintf("User with ID %d not found.", userID), http.StatusNotFound)
		return
	}

	// Attempt to delete the user
	err = app.Models.User.Delete(userID)
	if err != nil {
		helpers.RespondWithError(w, err, fmt.Sprintf("Failed to delete user with ID %d.", userID), http.StatusInternalServerError)
		return
	}

	// Log the deletion in the audit logs
	app.PushAuditToCache(
		*userData,
		"DELETE",
		"user",
		fmt.Sprintf("%d", userID),
		r,
		fmt.Sprintf("Deleted user with ID %d and email '%s'.", userID, user.Email),
	)

	// Respond with success
	response := map[string]interface{}{
		"message": fmt.Sprintf("User with ID %d successfully deleted.", userID),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		helpers.RespondWithError(w, err, "Failed to send response.", http.StatusInternalServerError)
	}
}

// Helper function to get client IP
func getClientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback to returning the whole field
	}
	return ip
}
