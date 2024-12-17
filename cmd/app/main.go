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

	// Start cron
	app.Cron.AddFunc("0,20,40 * * * * *", func() {
		app.Service.SyncRawLogs()
		app.Service.RegisterNewDevices()
		app.Service.SyncActivityLogs()
		app.Service.SyncDevices()
		app.Service.SyncNBIoTKeepaliveLogs()
		app.Service.SyncNBIoTSettingLogs()
		app.Service.SyncAuditLogs()
	})
	app.Cron.Start()

	// Create and start the the web server.
	httpServer := httpserver.NewServer(os.Getenv("HTTP_PORT"))
	httpServer.Start()
	defer httpServer.Shutdown()
}

// ---------------------------------------------------------------------

func initializeAppConfig() {
	time.Sleep(time.Millisecond * 10)

	// Load InfoLog and ErrorLog
	app.InfoLog = helpers.GetInfoLog()
	app.ErrorLog = helpers.GetErrorLog()
	app.FatalLog = helpers.GetFatalLog()

	var deviceAccessMode = os.Getenv("DEVICE_ACCESS_MODE")
	app.DeviceAccessMode = &deviceAccessMode

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
	app.Cache = cache.NewCache(redisPool, os.Getenv("REDIS_PREFIX"))

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

	app.InfoLog.Printf("App running in environment: %s\n", os.Getenv("GO_ENV"))
	return nil
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

	app.InfoLog.Println("Root user created or updated and cached successfully")
}

func initializeAppSettings() {
	// Check if the settings initialization is cached to avoid reinitialization
	isCached, err := app.Cache.Exists("app:settings")
	if err != nil {
		helpers.LogError(err, "Error checking cache for app settings")
	}
	if isCached {
		return // If settings are already cached, no need to reinitialize
	}

	// Encrypt the Google API key before storing it in the database
	googleApiKey, err := helpers.EncryptAES(os.Getenv("GOOGLE_API_KEY"), core.AES_SECRET_KEY)
	if err != nil {
		helpers.LogFatal(err, "Failed to encrypt Google API Key")
		return
	}

	// Prepare the settings data
	var settings = []models.Setting{
		{
			Key:         "google_api_key",
			Val:         googleApiKey,
			Description: "API key used for accessing Google services like Maps and Places.",
			AccessLevel: 0, // Root access level
			UpdatedBy:   0,
		},
		{
			Key:         "device_access_mode",
			Val:         os.Getenv("DEVICE_ACCESS_MODE"),
			Description: "Defines the access control mode for devices, determining whether they are managed via a blacklist or whitelist approach.",
			AccessLevel: 1, // Administrator access level
			UpdatedBy:   0,
		},
		{
			Key:         "default_latitude",
			Val:         os.Getenv("DEFAULT_LATITUDE"),
			Description: "Default latitude for map centering and initial device placement on the map.",
			AccessLevel: 1, // Administrator access level
			UpdatedBy:   0,
		},
		{
			Key:         "default_longitude",
			Val:         os.Getenv("DEFAULT_LONGITUDE"),
			Description: "Default longitude for map centering and initial device placement on the map.",
			AccessLevel: 1, // Administrator access level
			UpdatedBy:   0,
		},
		{
			Key:         "jwt_expiration_seconds",
			Val:         os.Getenv("JWT_EXPIRATION_TIME"),
			Description: "Duration in seconds for which a user's JSON Web Token (JWT) remains valid after login.",
			AccessLevel: 0, // Root access level
			UpdatedBy:   0,
		},
		{
			Key:         "redis_ttl_seconds",
			Val:         os.Getenv("REDIS_DEFAULT_TTL"),
			Description: "Default time-to-live (TTL) in seconds for items stored in the Redis cache, impacting how long user and device data are cached.",
			AccessLevel: 0, // Root access level
			UpdatedBy:   0,
		},
	}

	// Insert or update settings in the database
	for _, setting := range settings {
		_, err := setting.Upsert(&setting)
		if err != nil {
			helpers.LogFatal(err, "Failed to initialize application setting: "+setting.Key)
			continue // Optionally continue on error, depends on your error handling strategy
		}
	}
}
