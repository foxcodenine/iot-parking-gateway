package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"

	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/routes"
	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/core"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"

	"github.com/foxcodenine/iot-parking-gateway/internal/db"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/foxcodenine/iot-parking-gateway/internal/udp"
	"github.com/joho/godotenv"
)

var app core.App // Holds application-wide dependencies

func main() {
	// Initialize configuration and environment variables
	initializeAppConfig()

	// Initialize database connection
	initializeDatabase()
	defer app.DB.Close()

	// Initialize and set up handlers with app configuration
	initializeHandlers()

	// Start the UDP server in a goroutine
	udpServer := udp.NewServer(
		fmt.Sprintf(":%s", os.Getenv("UDP_PORT")),
		&app,
	)
	go func() {
		if err := udpServer.Start(); err != nil {
			app.ErrorLog.Fatalf("Failed to start UDP server: %v", err)
		}
	}()
	defer udpServer.Stop() // Ensure the UDP server is stopped on exit

	// Start the web server
	startServer()
}

// ---------------------------------------------------------------------

func initializeAppConfig() {
	// Load InfoLog and ErrorLog
	app.InfoLog = helpers.GetInfoLog()
	app.ErrorLog = helpers.GetErrorLog()

	// Load environment and configuration
	if err := loadEnv(); err != nil {
		app.ErrorLog.Fatalf("Error loading environment variables: %v\n", err)
	}
	app.AppURL = os.Getenv("APP_URL")
	app.HttpPort = os.Getenv("HTTP_PORT")

	// Register model types for gob encoding
	gob.Register(models.Device{})

	redisPool, err := cache.CreateRedisPool()

	if err != nil {
		helpers.LogError(err, "Failed to connect to Redis")
	} else {
		helpers.LogInfo("Successfully connected to Redis on :%s", os.Getenv("REDIS_PORT"))
	}

	app.Cache = cache.RedisCache{
		Conn:   redisPool,
		Prefix: os.Getenv("REDIS_PREFIX"),
	}
}

// ---------------------------------------------------------------------

func loadEnv() error {
	env := os.Getenv("GO_ENV")

	app.InfoLog.Printf("App running in environment: %s\n", os.Getenv("GO_ENV"))

	switch env {
	case "production":
		// return godotenv.Load("/app/.env") // Load production environment
		return nil
	default:
		return godotenv.Load(".env.development") // Load development environment
	}
}

// ---------------------------------------------------------------------

func initializeDatabase() {
	// Initialize database connection
	database, err := db.OpenDB()
	if err != nil {
		app.ErrorLog.Fatalf("Error connecting to the database: %v", err)
	}
	// Assign the database connection to the app configuration
	app.DB = database

	app.Models, err = models.New(database)
	if err != nil {
		app.ErrorLog.Fatalf("Error initializing models: %v", err)
	}
}

// ---------------------------------------------------------------------

func initializeHandlers() {
	// Initialize and set up the handlers with the application configuration
	handlersRepo := handlers.Initialize(&app)
	handlers.SetHandlerRepository(handlersRepo)
}

// ---------------------------------------------------------------------

func startServer() {
	// Set up the HTTP server configuration
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.HttpPort), // Bind the server to port
		Handler: routes.Routes(),                  // Set up the HTTP routes
	}

	app.InfoLog.Printf("HTTP server start on %s:%s\n", app.AppURL, app.HttpPort)

	// Start the server and handle any startup errors
	if err := srv.ListenAndServe(); err != nil {
		app.ErrorLog.Fatalf("Error starting the server: %v", err)
	}
}

// ---------------------------------------------------------------------
