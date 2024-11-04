package handlers

import (
	"github.com/foxcodenine/iot-parking-gateway/internal/core"
)

// Repository holds shared application configurations (App)
type Repository struct {
	App *core.App
}

// Initialize initializes the Repository with App and returns it
func Initialize(app *core.App) *Repository {
	return &Repository{App: app}
}

// repo is a global instance of Repository, allowing access to shared configurations
var repo *Repository
var app *core.App

// SetHandlerRepository sets the global repository (used globally if needed)
func SetHandlerRepository(r *Repository) {
	repo = r
	app = r.App
}

// GetRepo returns the global repository instance
func GetRepo() *Repository {
	return repo
}
