package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/apptypes"
	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/core"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

type SettingHandler struct {
}

func (h *SettingHandler) Update(w http.ResponseWriter, r *http.Request) {

	var rootLevelSettingsChange = false

	userData, err := app.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Authentication error.", http.StatusUnauthorized)
		return
	}

	// Check if the user has permission to update a user
	if userData.AccessLevel > 1 {
		http.Error(w, "You do not have the necessary permissions to perform this action.", http.StatusForbidden)
		return
	}

	// Define a map to hold the fields to update
	var updatedFields map[string]string

	// Decode the JSON body into the map for flexibility
	err = json.NewDecoder(r.Body).Decode(&updatedFields)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	admin_password, ok := updatedFields["admin_password"]
	if !ok {
		http.Error(w, "Missing required field: admin_password.", http.StatusForbidden)
		return
	}

	// Verify admin password
	adminUser, err := app.Models.User.FindUserByID(userData.UserID)
	if err != nil || !helpers.CheckPasswordHash(admin_password, adminUser.Password) {
		http.Error(w, "Authentication failed: incorrect admin credentials.", http.StatusForbidden)
		return
	}

	// -----------------------------------------------------------------

	google_api_key, ok := updatedFields["google_api_key"]
	if ok {

		if userData.AccessLevel > 0 {
			http.Error(w, "Permission denied.", http.StatusForbidden)
			return
		}

		// Validate that val is like a Google API key
		// Typical Google API keys are 39 characters long and consist of alphanumeric characters, dashes, and underscores
		apiKeyPattern := `^[A-Za-z0-9_-]{20,}$`
		matched, err := regexp.MatchString(apiKeyPattern, google_api_key)
		if err != nil || !matched {
			http.Error(w, "Invalid Google API key format.", http.StatusBadRequest)
			return
		}

		hashedGoogleApiKey, err := helpers.EncryptAES(google_api_key, core.AES_SECRET_KEY)
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to encrypt the Google API key.", http.StatusInternalServerError)
			return
		}

		val := map[string]interface{}{"val": hashedGoogleApiKey}

		// Proceed to update the setting
		_, err = app.Models.Setting.UpdateByKey("google_api_key", val)
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to update settings.", http.StatusInternalServerError)
			return
		}

		err = auditLogSettingUpdate(userData, r, "google_api_key", google_api_key)
		if err != nil {
			helpers.LogError(err, "Failed to create an audit log entry for updating the 'google_api_key' setting.")
		}
		rootLevelSettingsChange = true
	}

	// -----------------------------------------------------------------

	jwt_expiration_seconds, ok := updatedFields["jwt_expiration_seconds"]
	if ok {
		if userData.AccessLevel > 0 {
			http.Error(w, "Permission denied.", http.StatusForbidden)
			return
		}

		// Fetch the cached value
		cachedVal, err := cache.AppCache.HGet("app:settings", "jwt_expiration_seconds")
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to fetch cached setting value.", http.StatusInternalServerError)
			return
		}

		// Check if the value is the same
		if cachedVal.(string) != jwt_expiration_seconds {

			// Validate that the string represents an integer
			intVal, err := strconv.Atoi(jwt_expiration_seconds)
			if err != nil || intVal <= 0 {
				http.Error(w, "Invalid value for jwt_expiration_seconds. It must be a positive integer.", http.StatusBadRequest)
				return
			}

			val := map[string]interface{}{"val": jwt_expiration_seconds}

			// Proceed to update the setting
			_, err = app.Models.Setting.UpdateByKey("jwt_expiration_seconds", val)
			if err != nil {
				helpers.RespondWithError(w, err, "Failed to update settings.", http.StatusInternalServerError)
				return
			}

			err = auditLogSettingUpdate(userData, r, "jwt_expiration_seconds", jwt_expiration_seconds)
			if err != nil {
				helpers.LogError(err, "Failed to create an audit log entry for updating the 'jwt_expiration_seconds' setting.")
			}
			rootLevelSettingsChange = true
		}
	}

	// -----------------------------------------------------------------

	redis_ttl_seconds, ok := updatedFields["redis_ttl_seconds"]
	if ok {
		if userData.AccessLevel > 0 {
			http.Error(w, "Permission denied.", http.StatusForbidden)
			return
		}

		// Fetch the cached value
		cachedVal, err := cache.AppCache.HGet("app:settings", "redis_ttl_seconds")
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to fetch cached setting value.", http.StatusInternalServerError)
			return
		}

		// Check if the value is the same
		if cachedVal.(string) != redis_ttl_seconds {

			// Validate that the string represents an integer
			intVal, err := strconv.Atoi(redis_ttl_seconds)
			if err != nil || intVal <= 0 {
				http.Error(w, "Invalid value for redis_ttl_seconds. It must be a positive integer.", http.StatusBadRequest)
				return
			}

			val := map[string]interface{}{"val": redis_ttl_seconds}

			// Proceed to update the setting
			_, err = app.Models.Setting.UpdateByKey("redis_ttl_seconds", val)
			if err != nil {
				helpers.RespondWithError(w, err, "Failed to update settings.", http.StatusInternalServerError)
				return
			}

			err = auditLogSettingUpdate(userData, r, "redis_ttl_seconds", redis_ttl_seconds)
			if err != nil {
				helpers.LogError(err, "Failed to create an audit log entry for updating the 'redis_ttl_seconds' setting.")
			}
			rootLevelSettingsChange = true
		}
	}

	// -----------------------------------------------------------------

	device_access_mode, ok := updatedFields["device_access_mode"]
	if ok {

		if userData.AccessLevel > 0 {
			http.Error(w, "Permission denied.", http.StatusForbidden)
			return
		}

		// Fetch the cached value
		cachedVal, err := cache.AppCache.HGet("app:settings", "device_access_mode")
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to fetch cached setting value.", http.StatusInternalServerError)
			return
		}

		// Check if the value is the same
		if cachedVal.(string) != device_access_mode {

			// Validate that the string is either "black_list" or "white_list"
			if device_access_mode != "black_list" && device_access_mode != "white_list" {
				http.Error(w, "Invalid value for device_access_mode. It must be either 'black_list' or 'white_list'.", http.StatusBadRequest)
				return
			}

			val := map[string]interface{}{"val": device_access_mode}

			// Proceed to update the setting
			_, err = app.Models.Setting.UpdateByKey("device_access_mode", val)
			if err != nil {
				helpers.RespondWithError(w, err, "Failed to update settings.", http.StatusInternalServerError)
				return
			}

			err = auditLogSettingUpdate(userData, r, "device_access_mode", device_access_mode)
			if err != nil {
				helpers.LogError(err, "Failed to create an audit log entry for updating the 'device_access_mode' setting.")
			}
			rootLevelSettingsChange = true
		}
	}

	// -----------------------------------------------------------------

	loginPageTitle, ok := updatedFields["login_page_title"]
	if ok {
		if userData.AccessLevel > 0 {
			http.Error(w, "Permission denied.", http.StatusForbidden)
			return
		}

		// Fetch the cached value
		cachedVal, err := cache.AppCache.HGet("app:settings", "login_page_title")
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to fetch cached setting value.", http.StatusInternalServerError)
			return
		}

		// Check if the value is the same
		if cachedVal.(string) != loginPageTitle {
			// Validate that the string is not empty, has at least 5 characters, and does not contain potentially harmful tags
			if len(loginPageTitle) < 5 || regexp.MustCompile(`(?i)<script>|<iframe>|<object>|<embed>|<applet>`).MatchString(loginPageTitle) {
				http.Error(w, "Invalid value for login_page_title. It must be at least 5 characters long and free of any potentially harmful HTML tags.", http.StatusBadRequest)
				return
			}

			val := map[string]interface{}{"val": loginPageTitle}

			// Proceed to update the setting
			_, err = app.Models.Setting.UpdateByKey("login_page_title", val)
			if err != nil {
				helpers.RespondWithError(w, err, "Failed to update the 'login_page_title' setting.", http.StatusInternalServerError)
				return
			}

			err = auditLogSettingUpdate(userData, r, "login_page_title", loginPageTitle)
			if err != nil {
				helpers.LogError(err, "Failed to create an audit log entry for updating the 'login_page_title' setting.")
			}
			rootLevelSettingsChange = true
		}
	}

	// -----------------------------------------------------------------

	default_latitude, latOk := updatedFields["default_latitude"]
	if latOk {

		if userData.AccessLevel > 1 {
			http.Error(w, "Permission denied.", http.StatusForbidden)
			return
		}

		// Fetch the cached value
		cachedLat, err := cache.AppCache.HGet("app:settings", "default_latitude")
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to fetch cached setting value.", http.StatusInternalServerError)
			return
		}

		// Check if the value is the same
		if cachedLat.(string) != default_latitude {

			// Validate that the string represents a float
			_, err := strconv.ParseFloat(default_latitude, 64)
			if err != nil {
				http.Error(w, "Invalid value for default_latitude. It must be a valid floating-point number.", http.StatusBadRequest)
				return
			}

			val := map[string]interface{}{"val": default_latitude}

			// Proceed to update the setting
			_, err = app.Models.Setting.UpdateByKey("default_latitude", val)
			if err != nil {
				helpers.RespondWithError(w, err, "Failed to update settings.", http.StatusInternalServerError)
				return
			}

			err = auditLogSettingUpdate(userData, r, "default_latitude", default_latitude)
			if err != nil {
				helpers.LogError(err, "Failed to create an audit log entry for updating the 'default_latitude' setting.")
			}
		}
	}

	// -----------------------------------------------------------------

	default_longitude, lonOk := updatedFields["default_longitude"]
	if lonOk {

		if userData.AccessLevel > 1 {
			http.Error(w, "Permission denied.", http.StatusForbidden)
			return
		}

		// Fetch the cached value
		cachedLon, err := cache.AppCache.HGet("app:settings", "default_longitude")
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to fetch cached setting value.", http.StatusInternalServerError)
			return
		}

		// Check if the value is the same
		if cachedLon.(string) != default_longitude {

			// Validate that the string represents a float
			_, err := strconv.ParseFloat(default_longitude, 64)
			if err != nil {
				http.Error(w, "Invalid value for default_longitude. It must be a valid floating-point number.", http.StatusBadRequest)
				return
			}

			val := map[string]interface{}{"val": default_longitude}

			// Proceed to update the setting
			_, err = app.Models.Setting.UpdateByKey("default_longitude", val)
			if err != nil {
				helpers.RespondWithError(w, err, "Failed to update settings.", http.StatusInternalServerError)
				return
			}
			err = auditLogSettingUpdate(userData, r, "default_longitude", default_longitude)
			if err != nil {
				helpers.LogError(err, "Failed to create an audit log entry for updating the 'default_longitude' setting.")
			}
		}
	}

	// -----------------------------------------------------------------

	google_map_id, ok := updatedFields["google_map_id"]
	if ok {
		if userData.AccessLevel > 0 {
			http.Error(w, "Permission denied.", http.StatusForbidden)
			return
		}

		// Fetch the cached value
		cachedVal, err := cache.AppCache.HGet("app:settings", "google_map_id")
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to fetch cached setting value.", http.StatusInternalServerError)
			return
		}

		// Check if the value is the same
		if cachedVal.(string) != google_map_id {
			// Validate that the string represents a valid Google Map ID pattern
			mapIDPattern := `^[A-Za-z0-9_-]{10,}$` // Placeholder pattern; adjust as per actual requirements
			matched, err := regexp.MatchString(mapIDPattern, google_map_id)
			if err != nil || !matched {
				http.Error(w, "Invalid Google Map ID format.", http.StatusBadRequest)
				return
			}

			val := map[string]interface{}{"val": google_map_id}

			// Proceed to update the setting
			_, err = app.Models.Setting.UpdateByKey("google_map_id", val)
			if err != nil {
				helpers.RespondWithError(w, err, "Failed to update the 'google_map_id' setting.", http.StatusInternalServerError)
				return
			}

			err = auditLogSettingUpdate(userData, r, "google_map_id", google_map_id)
			if err != nil {
				helpers.LogError(err, "Failed to create an audit log entry for updating the 'google_map_id' setting.")
			}
			rootLevelSettingsChange = true
		}
	}

	// -----------------------------------------------------------------

	cors_allowed_origins, ok := updatedFields["cors_allowed_origins"]
	if ok {
		if userData.AccessLevel > 0 {
			http.Error(w, "Permission denied.", http.StatusForbidden)
			return
		}

		// Fetch the cached value
		cachedVal, err := cache.AppCache.HGet("app:settings", "cors_allowed_origins")
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to fetch cached setting value.", http.StatusInternalServerError)
			return
		}

		// Check if the value is the same
		if cachedVal.(string) != cors_allowed_origins {
			// Validate the value as a comma-separated list of domains or '*'
			corsPattern := `^(\*,|https?://[\w.-]+(:\d+)?(,|$))*$`
			matched, err := regexp.MatchString(corsPattern, cors_allowed_origins)
			if err != nil || !matched {
				http.Error(w, "Invalid CORS allowed origins format. Use '*' or a comma-separated list of valid URLs.", http.StatusBadRequest)
				return
			}

			val := map[string]interface{}{"val": cors_allowed_origins}

			// Proceed to update the setting
			_, err = app.Models.Setting.UpdateByKey("cors_allowed_origins", val)
			if err != nil {
				helpers.RespondWithError(w, err, "Failed to update the 'cors_allowed_origins' setting.", http.StatusInternalServerError)
				return
			}

			err = auditLogSettingUpdate(userData, r, "cors_allowed_origins", cors_allowed_origins)
			if err != nil {
				helpers.LogError(err, "Failed to create an audit log entry for updating the 'cors_allowed_origins' setting.")
			}
			rootLevelSettingsChange = true
		}
	}

	// -----------------------------------------------------------------

	initial_parking_check_date, ok := updatedFields["initial_parking_check_date"]
	if ok {
		if userData.AccessLevel > 1 {
			http.Error(w, "Permission denied.", http.StatusForbidden)
			return
		}

		// Fetch the cached value
		cachedVal, err := cache.AppCache.HGet("app:settings", "initial_parking_check_date")
		if err != nil {
			helpers.RespondWithError(w, err, "Failed to fetch cached setting value.", http.StatusInternalServerError)
			return
		}

		// Check if the value is the same
		if cachedVal.(string) != initial_parking_check_date {
			// Validate that the value is in ISO 8601 format (e.g., 2014-12-21T15:35:24Z)
			datePattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`
			matched, err := regexp.MatchString(datePattern, initial_parking_check_date)
			if err != nil || !matched {
				http.Error(w, "Invalid initial parking check date format. It must be in ISO 8601 format (e.g., 2014-12-21T15:35:24Z).", http.StatusBadRequest)
				return
			}

			val := map[string]interface{}{"val": initial_parking_check_date}

			// Proceed to update the setting
			_, err = app.Models.Setting.UpdateByKey("initial_parking_check_date", val)
			if err != nil {
				helpers.RespondWithError(w, err, "Failed to update the 'initial_parking_check_date' setting.", http.StatusInternalServerError)
				return
			}

			err = auditLogSettingUpdate(userData, r, "initial_parking_check_date", initial_parking_check_date)
			if err != nil {
				helpers.LogError(err, "Failed to create an audit log entry for updating the 'initial_parking_check_date' setting.")
			}
			rootLevelSettingsChange = true
		}
	}

	// -----------------------------------------------------------------

	var message = "Settings updated. Changes will take effect upon next login."

	if rootLevelSettingsChange {
		message = "Settings updated. All users will be logged out to reload new settings."

		// Generate a new JWT secret key
		secretKey, err := helpers.GenerateJWTSecretKey(44)
		if err != nil {
			helpers.LogError(err, "Failed to generate a new JWT secret key.")
		}

		// Update the environment variable
		err = os.Setenv("JWT_SECRET_KEY", secretKey)
		if err != nil {
			helpers.LogError(err, "Failed to set new JWT secret key in environment.")
		}
	}
	response := map[string]interface{}{
		"message": message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		helpers.RespondWithError(w, err, "Failed to encode response.", http.StatusInternalServerError)
	}
}

func auditLogSettingUpdate(userData *apptypes.UserClaims, r *http.Request, key, newValue string) error {
	// Create the audit log entry
	auditLogEntry := models.AuditLog{
		UserID:      userData.UserID,
		Email:       userData.Email,
		AccessLevel: userData.AccessLevel,
		HappenedAt:  time.Now().UTC(),
		Action:      "UPDATE",
		Entity:      "SETTING",
		URL:         r.URL.Path,
		IPAddress:   getClientIP(r),
		Details:     fmt.Sprintf("User with ID %d and email '%s' updated '%s' to '%s'.", userData.UserID, userData.Email, key, newValue),
	}

	// Push the audit log entry to the cache
	return app.Cache.RPush("logs:audit-logs", auditLogEntry)
}
