package handlers

import "github.com/foxcodenine/iot-parking-gateway/internal/config"

// Repository holds shared application configurations (AppConfig)
type Repository struct {
	App *config.App
}

// Initialize initializes the Repository with AppConfig and returns it
func Initialize(app *config.App) *Repository {
	return &Repository{App: app}
}

// repo is a global instance of Repository, allowing access to shared configurations
var repo *Repository
var app *config.App

// SetHandlerRepository sets the global repository (used globally if needed)
func SetHandlerRepository(r *Repository) {
	repo = r
	app = r.App
}

// GetRepo returns the global repository instance
func GetRepo() *Repository {
	return repo
}
