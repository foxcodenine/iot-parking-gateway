package main

import (
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

var app config.App
var appUrl string
var webPort string

func main() {
	fmt.Printf("App running in - %s\n", os.Getenv("GO_ENV"))
	fmt.Printf("Starting web server - %s:%s/\n", appUrl, webPort)

	db, err := db.OpenDB()

	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}

	app.DB = db

	models, err := models.New(db)

	if err != nil {
		log.Fatalf("Error creating models: %v", err)
	}

	app.Models = models

	models.Device.Create("0123456789ABCDE")

	// Initialize and set the handlers repository
	handlersRepo := handlers.Initialize(&app)
	handlers.SetHandlerRepository(handlersRepo)

	// Start the server
	serv()
}

func init() {
	// Load environment variables based on environment (production, development, etc.)
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Set application URL and port from environment variables
	appUrl = os.Getenv("APP_URL")
	webPort = os.Getenv("APP_WEB_PORT")
}

// loadEnv loads the appropriate environment file based on GO_ENV
func loadEnv() error {
	env := os.Getenv("GO_ENV")

	switch env {
	case "production":
		return godotenv.Load("/app/.env")
	default:
		return godotenv.Load(".env.development")
	}
}

// serv initializes and starts the HTTP server with routes and configuration
func serv() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routes.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
