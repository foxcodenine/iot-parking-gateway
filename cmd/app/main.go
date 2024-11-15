package main

import (
	"encoding/gob"
	"fmt"

	"os"

	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/handlers"
	"github.com/foxcodenine/iot-parking-gateway/internal/mq"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/core"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/httpserver"
	"github.com/foxcodenine/iot-parking-gateway/internal/services"
	"github.com/robfig/cron/v3"

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
	go app.UdpServer.Start()
	defer app.UdpServer.Stop()

	// Start cron
	app.Cron.AddFunc("* * * * *", func() {
		app.Service.TransferRawLogsFromRedisToPostgres()
		app.Service.TransferActivityLogsFromRedisToPostgres()
	})
	app.Cron.Start()

	// Setup RabbitMQ Producer
	rabbitConfig := mq.SetupRabbitMQConfig()
	rabbitProducer := mq.NewRabbitMQProducer(rabbitConfig)
	defer rabbitProducer.Close() // Ensure to close the connection on application shutdown

	// Start the message producer routine
	go rabbitProducer.Run()

	// Start the web server
	httpServer := httpserver.NewServer(os.Getenv("HTTP_PORT"))
	httpServer.Start()
	defer httpServer.Shutdown()
}

// ---------------------------------------------------------------------

func initializeAppConfig() {
	time.Sleep(time.Millisecond * 10)
	fmt.Println("")

	// Load InfoLog and ErrorLog
	app.InfoLog = helpers.GetInfoLog()
	app.ErrorLog = helpers.GetErrorLog()
	app.FatalLog = helpers.GetFatalLog()

	// Load environment and configuration
	if err := loadEnv(); err != nil {
		app.ErrorLog.Fatalf("Error loading environment variables: %v\n", err)
	}
	app.AppURL = os.Getenv("APP_URL")
	app.HttpPort = os.Getenv("HTTP_PORT")

	// Register model types for gob encoding
	gob.Register(models.Device{})

	// Initialize a Redis connection pool
	redisPool, err := cache.CreateRedisPool()
	if err != nil {
		app.ErrorLog.Fatalf("Failed to connect to Redis: %v\n", err)
	} else {
		app.InfoLog.Printf("Successfully connected to Redis on :%s", os.Getenv("REDIS_PORT"))
	}

	// Assign Redis cache instance to the app configuration
	app.Cache = &cache.RedisCache{
		Conn:   redisPool,
		Prefix: os.Getenv("REDIS_PREFIX"), // Use a prefix for cache keys, if provided
	}

	app.Service = services.NewService(
		app.Models,
		app.Cache,
		app.InfoLog,
		app.ErrorLog,
	)

	app.UdpServer = udp.NewServer(
		fmt.Sprintf(":%s", os.Getenv("UDP_PORT")),
		app.Cache,
		app.Service,
	)

	// Initialize and assign a cron scheduler instance to the app
	app.Cron = cron.New()
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
