package handlers

import "github.com/foxcodenine/iot-parking-gateway/internal/appconfig"

// Repository holds shared application configurations (AppConfig)
type Repository struct {
	App *appconfig.AppConfig
}

// Initialize initializes the Repository with AppConfig and returns it
func Initialize(app *appconfig.AppConfig) *Repository {
	return &Repository{App: app}
}

// Repo is a global instance of Repository, allowing access to shared configurations
var Repo *Repository

// SetHandlerRepository sets the global repository (used globally if needed)
func SetHandlerRepository(r *Repository) {
	Repo = r
}

// GetRepo returns the global repository instance
func GetRepo() *Repository {
	return Repo
}
