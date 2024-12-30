package services

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/go-faker/faker/v4"
)

type Service struct {
	models   models.Models
	cache    *cache.RedisCache
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewService(m models.Models, rc *cache.RedisCache, il, el *log.Logger) *Service {
	return &Service{
		models:   m,
		cache:    rc,
		infoLog:  il,
		errorLog: el,
	}
}

// RegisterNewDevices retrieves device IDs from Redis, creates new device records in the database, and logs the outcome.
func (s *Service) RegisterNewDevices() {
	// Retrieve device IDs from Redis and delete the set afterwards.
	deviceEntries, err := s.cache.SMembersDel("to-register-devices")
	if err != nil {
		s.errorLog.Printf("Failed to retrieve device IDs from Redis: %v", err)
		return
	}

	// Process each device entry to create new device records.
	for _, entry := range deviceEntries {

		// Convert interface to string, ensuring it represents a device ID correctly.
		deviceInfo, ok := entry.(string)
		if !ok {
			s.errorLog.Printf("Invalid device ID format: expected string but got %T", entry)
			continue
		}

		// Split the string into parts based on spaces
		parts := strings.Split(deviceInfo, " ")
		if len(parts) < 3 {
			s.errorLog.Printf("Device information string does not have enough parts: %s\n", deviceInfo)
			return // or continue if in a loop
		}

		// Assign parts to specific variables
		networkType := parts[0]
		deviceID := parts[1]
		firmwareVersionStr := parts[2]

		// Convert firmwareVersionStr to a float
		firmwareVersion, err := strconv.ParseFloat(firmwareVersionStr, 64) // 64 specifies the precision
		if err != nil {
			s.errorLog.Printf("Failed to convert firmware version to float: %s\n", err)
			return // or handle the error appropriately
		}

		latitude, err := strconv.ParseFloat(os.Getenv("DEFAULT_LATITUDE"), 64)
		if err != nil {
			latitude = 0
		}
		longitude, err := strconv.ParseFloat(os.Getenv("DEFAULT_LONGITUDE"), 64)
		if err != nil {
			longitude = 0
		}

		name := faker.Username()
		// Define a new device model instance.
		newDevice := models.Device{
			DeviceID:        deviceID,
			Name:            "__" + name,
			NetworkType:     networkType,
			Latitude:        latitude,
			Longitude:       longitude,
			FirmwareVersion: firmwareVersion,
		}

		// Attempt to create a new device record in the database.
		_, err = s.models.Device.Create(&newDevice)
		if err != nil {
			s.errorLog.Printf("Failed to create a new device record for ID %s: %v", deviceID, err)
			continue // Proceed to the next ID instead of stopping.
		}

		// If the network type is 'nb', attempts to create a new device settings record in the database.
		if newDevice.NetworkType == "NB-IoT" {
			newDeviceSetting := models.NbiotDeviceSettings{
				DeviceID:    deviceID,
				NetworkType: networkType,
			}
			_, err = s.models.NbiotDeviceSettings.Create(&newDeviceSetting)
			if err != nil {
				s.errorLog.Printf("Failed to create a new device settings record for ID %s: %v", deviceID, err)
				continue // Proceed to the next ID instead of stopping.
			}
		}

	}
	if len(deviceEntries) == 1 {
		s.infoLog.Println("Successfully added 1 device to PostgreSQL")
	} else if len(deviceEntries) > 1 {
		s.infoLog.Printf("Successfully added %d devices to PostgreSQL", len(deviceEntries))
	}
}

// PopulateDeviceBloomFilter retrieves all devices from the database and populates a Bloom Filter with their network type and device IDs.
func (s *Service) PopulateDeviceBloomFilter() {
	// Fetch all devices from the database.
	deviceInstance := models.Device{}
	devices, err := deviceInstance.GetAllIncludingDeleted()
	if err != nil {
		s.errorLog.Printf("Error retrieving devices from Postgres: %v", err)
		return
	}

	// Initialize a slice to store composite keys for the Bloom Filter.
	var keyStrings []string

	// Iterate over each device to create a composite key of network type and device ID.
	for _, d := range devices {
		key := fmt.Sprintf("%s %s", d.NetworkType, d.DeviceID)
		keyStrings = append(keyStrings, key)
	}

	// Add all composite keys to the Bloom Filter named "registered-devices".
	s.cache.AddItemsToBloomFilter("registered-devices", keyStrings)
}

func (s *Service) PopulateDeviceCache() {
	// Step 1: Clear the device cache
	if err := s.cache.DeleteAllDevices(); err != nil {
		helpers.LogError(err, "Failed to clear device cache")
		return
	}
	helpers.LogInfo("Device cache cleared successfully.")

	// Step 2: Retrieve all devices from the database
	allDevices, err := s.models.Device.GetAll()
	if err != nil {
		helpers.LogError(err, "Failed to retrieve devices from database")
		return
	}
	helpers.LogInfo("Retrieved %d devices from database.", len(allDevices))

	// Step 3: Convert the slice of devices to []map[string]any
	deviceMaps, err := helpers.StructSliceToMapSlice(allDevices)
	if err != nil {
		helpers.LogError(err, "Failed to convert devices to map slice")
		return
	}

	// Step 4: Save all devices to the cache
	if err := s.cache.SaveMultipleDevices(deviceMaps); err != nil {
		helpers.LogError(err, "Failed to save devices to cache")
		return
	}
	helpers.LogInfo("Device cache populated successfully with %d devices.", len(allDevices))
}
