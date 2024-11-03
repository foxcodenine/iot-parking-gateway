package handlers

import "github.com/foxcodenine/iot-parking-gateway/internal/config"

// Repository holds shared application configurations (AppConfig)
type Repository struct {
	App *config.AppConfig
}

// Initialize initializes the Repository with AppConfig and returns it
func Initialize(app *config.AppConfig) *Repository {
	return &Repository{App: app}
}

// repo is a global instance of Repository, allowing access to shared configurations
var repo *Repository
var app *config.AppConfig

// SetHandlerRepository sets the global repository (used globally if needed)
func SetHandlerRepository(r *Repository) {
	repo = r
	app = r.App
}

// GetRepo returns the global repository instance
func GetRepo() *Repository {
	return repo
}
