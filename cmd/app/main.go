package main

import (
	"encoding/gob"
	"fmt"
	"log"
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
	initializeAppSettings()

	// Initialize and set up handlers with app configuration
	initializeHandlers()

	// Start the message producer routine
	go app.MQProducer.Run()
	defer app.MQProducer.Close() // Ensure to close the connection on application shutdown

	// Start the UDP server in a goroutine
	go app.UdpServer.Start()
	defer app.UdpServer.Stop()

	// // Initialize and Populate a Bloom Filter in Redis to efficiently check the existence of device IDs.
	app.Cache.CreateBloomFilter("registered-devices", 0.00001, 100000)
	app.Service.PopulateDeviceBloomFilter()
	app.Service.PopulateDeviceCache()

	// Start cron
	app.Cron.AddFunc("0,20,40 * * * * *", func() {
		app.Service.SyncRawLogs()
		app.Service.RegisterNewDevices()

		time.Sleep(1 * time.Second)

		app.Service.SyncDevices()
		app.Service.SyncActivityLogs()
		app.Service.SyncDevicesKeepaliveAt()
		app.Service.SyncDevicesSettingsAt()

		app.Service.SyncLoraKeepaliveLogs()
		app.Service.SyncLoraSettingLogs()

		app.Service.SyncNBIoTKeepaliveLogs()
		app.Service.SyncNBIoTSettingLogs()

		app.Service.SyncSigfoxKeepaliveLogs()
		app.Service.SyncSigfoxSettingLogs()

		app.Service.SyncAuditLogs()
	})
	app.Cron.Start()

	// Create and start the the web server.
	httpServer := httpserver.NewHttpServer(os.Getenv("HTTP_PORT"))
	httpServer.Start()
	defer httpServer.Shutdown()

	app.UdpServer.SocketIO = httpServer.SocketServer
	app.SocketIO = httpServer.SocketServer
}

// ---------------------------------------------------------------------

func initializeAppConfig() {
	time.Sleep(time.Millisecond * 10)

	var deviceAccessMode = os.Getenv("DEVICE_ACCESS_MODE")
	app.DeviceAccessMode = &deviceAccessMode

	// Load environment and configuration
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading environment variables \n%v", err)
	}
	helpers.ConfigLogger()
	helpers.LogInfo("App running in environment: %s", os.Getenv("GO_ENV"))

	app.AppURL = os.Getenv("APP_URL")
	app.HttpPort = os.Getenv("HTTP_PORT")

	// Register model types for gob encoding
	gob.Register(models.Device{})

	// Initialize a Redis connection pool
	redisPool, err := cache.CreateRedisPool()
	if err != nil {
		helpers.LogFatal(err, "Failed to connect to Redis")
	} else {
		helpers.LogInfo("Successfully connected to Redis on :%s", os.Getenv("REDIS_PORT"))
	}

	// Assign Redis cache instance to the app configuration
	app.Cache = cache.NewCache(redisPool, os.Getenv("REDIS_PREFIX"))

	// Initialize the service layer that handles business logic.
	app.Service = services.NewService(
		app.Models,
		app.Cache,
	)

	// Setup RabbitMQ Producer
	rabbitConfig := mq.SetupRabbitMQConfig()
	app.MQProducer = mq.NewRabbitMQProducer(rabbitConfig)

	// Set up the UDP server
	app.UdpServer = udp.NewUDPServer(
		fmt.Sprintf(":%s", os.Getenv("UDP_PORT")),
		app.MQProducer,
		app.Cache,
		app.Service,
		app.DeviceAccessMode,
	)

	// Initialize and assign a cron scheduler instance to the app
	app.Cron = cron.New(cron.WithSeconds())
}

// ---------------------------------------------------------------------

func loadEnv() error {
	env := os.Getenv("GO_ENV")

	switch env {
	case "production":
	default:
		godotenv.Load(".env.development") // Load development environment
	}

	// Generate a new JWT secret key
	secretKey, err := helpers.GenerateJWTSecretKey(44)
	if err != nil {
		return fmt.Errorf("failed to generate a new JWT secret key. \n%w", err)
	}

	// Update the environment variable
	err = os.Setenv("JWT_SECRET_KEY", secretKey)
	if err != nil {
		return fmt.Errorf("failed to set new JWT secret key in environment. \n%w", err)
	}
	return nil
}

// ---------------------------------------------------------------------

func initializeDatabase() {
	// Initialize database connection
	database, err := db.OpenDB()
	if err != nil {
		helpers.LogFatal(err, "Error connecting to the database")
	}
	// Assign the database connection to the app configuration
	app.DB = database

	app.Models, err = models.New(database)
	if err != nil {
		helpers.LogFatal(err, "Error initializing models")
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
	// Get the current root user details or a new user if not existing
	rootUser, err := app.Models.User.GetRootUser()
	if err != nil {
		helpers.LogError(err, "Failed to retrieve root user")
		return
	}

	if rootUser == nil {

		// If rootUser does not exist, initialize a new User struct
		rootUser = &models.User{
			Email:       os.Getenv("APP_ROOT_EMAIL"),
			Password:    os.Getenv("APP_ROOT_PASSWORD"),
			AccessLevel: 0, // Root access level
			Enabled:     true,
		}
		// Attempt to create the root user
		if _, err := rootUser.Create(); err != nil {
			helpers.LogError(err, "Failed to create root user")
			return
		}
	} else {
		// Update the existing root user details
		rootUser.Email = os.Getenv("APP_ROOT_EMAIL")
		rootUser.Password = os.Getenv("APP_ROOT_PASSWORD")
		rootUser.AccessLevel = 0
		rootUser.Enabled = true

		// Attempt to update the root user
		if _, err := rootUser.Update(true); err != nil {
			helpers.LogError(err, "Failed to update root user")
			return
		}
	}

	helpers.LogInfo("Root user created or updated and cached successfully")
}
