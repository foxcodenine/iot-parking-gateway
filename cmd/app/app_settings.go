package main

import (
	"os"

	"github.com/foxcodenine/iot-parking-gateway/internal/core"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

func initializeAppSettings() {
	// Check if the settings initialization is cached to avoid reinitialization
	isCached, err := app.Cache.Exists("app:settings")
	if err != nil {
		helpers.LogError(err, "Error checking cache for app settings")
	}
	if isCached {
		return // If settings are already cached, no need to reinitialize
	}

	// Encrypt the Google API key before storing it in the database
	googleApiKey, err := helpers.EncryptAES(os.Getenv("GOOGLE_API_KEY"), core.AES_SECRET_KEY)
	if err != nil {
		helpers.LogFatal(err, "Failed to encrypt Google API Key")
		return
	}

	// Prepare the settings data
	var settings = []models.Setting{
		{
			Key:         "google_api_key",
			Val:         googleApiKey,
			Description: "API key used for accessing Google services like Maps and Places.",
			AccessLevel: 0, // Root access level
			UpdatedBy:   0,
		},
		{
			Key:         "google_map_id",
			Val:         os.Getenv("GOOGLE_MAP_ID"),
			Description: "The Google Map ID used to customize and embed Google Maps in the application.",
			AccessLevel: 0,
			UpdatedBy:   0,
		},
		{
			Key:         "jwt_expiration_seconds",
			Val:         os.Getenv("JWT_EXPIRATION_TIME"),
			Description: "Duration in seconds for which a user's JSON Web Token (JWT) remains valid after login.",
			AccessLevel: 0, // Root access level
			UpdatedBy:   0,
		},
		{
			Key:         "redis_ttl_seconds",
			Val:         os.Getenv("REDIS_DEFAULT_TTL"),
			Description: "Default time-to-live (TTL) in seconds for items stored in the Redis cache.",
			AccessLevel: 0, // Root access level
			UpdatedBy:   0,
		},
		{
			Key:         "device_access_mode",
			Val:         os.Getenv("DEVICE_ACCESS_MODE"),
			Description: "Defines the access control mode for devices, determining whether they are managed via a blacklist or whitelist approach.",
			AccessLevel: 0,
			UpdatedBy:   0,
		},
		{
			Key:         "initial_parking_check_date",
			Val:         "2014-12-21T15:35:24Z",
			Description: "The reference date for checking parking events. Devices with no events after this date are considered newly installed or inactive, and their status is marked as unknown.",
			AccessLevel: 0,
			UpdatedBy:   0,
		},
		{
			Key:         "cors_allowed_origins",
			Val:         "*,http://localhost:5173,http://127.0.0.1:5173",
			Description: "Specifies the domains that are permitted to access the API, including development hosts. Use '*' to allow all or specify domains individually, separated by a comma.",
			AccessLevel: 0,
			UpdatedBy:   0,
		},
		{
			Key:         "default_latitude",
			Val:         os.Getenv("DEFAULT_LATITUDE"),
			Description: "Default latitude for map centering and initial device placement on the map.",
			AccessLevel: 1,
			UpdatedBy:   0,
		},
		{
			Key:         "default_longitude",
			Val:         os.Getenv("DEFAULT_LONGITUDE"),
			Description: "Default longitude for map centering and initial device placement on the map.",
			AccessLevel: 1,
			UpdatedBy:   0,
		},
		{
			Key:         "login_page_title",
			Val:         "Welcome to <b>IoTrack</b> Pro",
			Description: "The HTML-formatted title text displayed on the login page of the IoTrack Pro application.",
			AccessLevel: 0, // Root access level required for changes
			UpdatedBy:   0, // Updated by the system or root user by default
		},
	}

	// Insert or update settings in the database
	for _, setting := range settings {
		_, err := setting.Upsert(&setting)
		if err != nil {
			helpers.LogFatal(err, "Failed to initialize application setting: "+setting.Key)
			continue // Optionally continue on error, depends on your error handling strategy
		}
	}
}
