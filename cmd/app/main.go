package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/routes"
	"github.com/foxcodenine/iot-parking-gateway/internal/config"
	"github.com/foxcodenine/iot-parking-gateway/internal/db"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/joho/godotenv"
)

var app config.App // Holds application-wide dependencies
var appUrl string  // Application URL loaded from environment variables
var webPort string // Application port loaded from environment variables

func main() {

	// Display environment and server details
	app.InfoLog.Printf("App running in environment: %s\n", os.Getenv("GO_ENV"))
	app.InfoLog.Printf("Starting web server at %s:%s/\n", appUrl, webPort)

	// Initialize database connection
	database, err := db.OpenDB()
	if err != nil {
		app.ErrorLog.Fatalf("Error connecting to the database: %v", err)
	}

	// Ensure the database connection is closed on exit
	defer database.Close()

	// Assign the database connection to the app configuration
	app.DB = database

	// Initialize models and add them to the app configuration
	appModels, err := models.New(database)
	if err != nil {
		app.ErrorLog.Fatalf("Error creating models: %v", err)
	}
	app.Models = appModels

	// Initialize and set up the handlers with the application configuration
	handlersRepo := handlers.Initialize(&app)
	handlers.SetHandlerRepository(handlersRepo)

	// Start the web server
	startServer()
}

// ---------------------------------------------------------------------

func init() {
	// Initialize InfoLog and ErrorLog
	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Load environment variables based on the GO_ENV variable
	if err := loadEnv(); err != nil {
		app.ErrorLog.Fatalf("Error loading environment variables: %v", err)
	}

	// Set application URL and port from environment variables
	appUrl = os.Getenv("APP_URL")
	webPort = os.Getenv("APP_WEB_PORT")

	// Register model types for serialization
	gob.Register(models.Device{})
}

// ---------------------------------------------------------------------

// loadEnv loads environment variables from the appropriate file based on GO_ENV
func loadEnv() error {
	env := os.Getenv("GO_ENV")
	switch env {
	case "production":
		return godotenv.Load("/app/.env") // Load production environment
	default:
		return godotenv.Load(".env.development") // Load development environment
	}
}

// startServer initializes and starts the HTTP server with the configured routes
func startServer() {
	// Set up the HTTP server configuration
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort), // Bind the server to port
		Handler: routes.Routes(),             // Set up the HTTP routes
	}

	// Start the server and handle any startup errors
	if err := srv.ListenAndServe(); err != nil {
		app.ErrorLog.Fatalf("Error starting the server: %v", err)
	}
}

// ---------------------------------------------------------------------
