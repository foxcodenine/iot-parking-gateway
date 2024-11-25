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
	initializeRootUser()

	// Initialize and set up handlers with app configuration
	initializeHandlers()

	// Start the message producer routine
	go app.MQProducer.Run()
	defer app.MQProducer.Close() // Ensure to close the connection on application shutdown

	// Start the UDP server in a goroutine
	go app.UdpServer.Start()
	defer app.UdpServer.Stop()

	// // Initialize and Populate a Bloom Filter in Redis to efficiently check the existence of device IDs.
	app.Cache.CreateBloomFilter("device-id", 0.00001, 100000)
	app.Service.PopulateDeviceBloomFilter()

	// Start cron
	app.Cron.AddFunc("0,20,40 * * * * *", func() {
		app.Service.SyncRawLogs()
		app.Service.RegisterNewDevices()
		app.Service.SyncActivityLogsAndDevices()
		app.Service.SyncNBIoTKeepaliveLogs()
		app.Service.SyncNBIoTSettingLogs()
	})
	app.Cron.Start()

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

	// Initialize the service layer that handles business logic.
	app.Service = services.NewService(
		app.Models,
		app.Cache,
		app.InfoLog,
		app.ErrorLog,
	)

	// Setup RabbitMQ Producer
	rabbitConfig := mq.SetupRabbitMQConfig()
	app.MQProducer = mq.NewRabbitMQProducer(rabbitConfig)

	// Set up the UDP server
	app.UdpServer = udp.NewServer(
		fmt.Sprintf(":%s", os.Getenv("UDP_PORT")),
		app.MQProducer,
		app.Cache,
		app.Service,
	)

	// Initialize and assign a cron scheduler instance to the app
	app.Cron = cron.New(cron.WithSeconds())
}

// ---------------------------------------------------------------------

func loadEnv() error {
	env := os.Getenv("GO_ENV")
	app.InfoLog.Printf("App running in environment: %s\n", os.Getenv("GO_ENV"))

	switch env {
	case "production":
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

func initializeRootUser() {
	// Construct the cache key for the root user
	rootUserCacheKey := "initialized-root-user:" + os.Getenv("APP_ROOT_EMAIL")

	// Check if the root user initialization is cached
	isCached, _ := app.Cache.Exists(rootUserCacheKey)
	if isCached {
		return
	}

	// Check if the root user already exists in the database
	existingUser, _ := app.Models.User.FindUserByEmail(os.Getenv("APP_ROOT_EMAIL"))
	if existingUser != nil {
		// Mark the root user as initialized in the cache
		_ = app.Cache.Set(rootUserCacheKey, "true", -1) // -1 indicates no expiration
		return
	}

	// Create the root user with default credentials from environment variables
	rootUser := models.User{
		Email:       os.Getenv("APP_ROOT_EMAIL"),
		Password:    os.Getenv("APP_ROOT_PASSWORD"),
		AccessLevel: 0, // 0 represents the highest level of access
		Enabled:     true,
	}

	// Attempt to create the root user
	if _, err := rootUser.Create(); err != nil {
		fmt.Printf("Failed to create root user: %v\n", err)
		return
	}

	// Cache the root user initialization to prevent duplicate runs
	_ = app.Cache.Set(rootUserCacheKey, "true", -1)

	app.InfoLog.Println("Root user created successfully")
}
