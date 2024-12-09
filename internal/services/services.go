package services

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
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

		// Extract network type and device ID from the string.
		spaceIndex := strings.Index(deviceInfo, " ")
		if spaceIndex == -1 {
			s.errorLog.Printf("Device information string does not contain a space: %s", deviceInfo)
			continue
		}

		networkType := deviceInfo[:spaceIndex]
		deviceID := deviceInfo[spaceIndex+1:]

		latitude, err := strconv.ParseFloat(os.Getenv("DEFAULT_LATITUDE"), 64)
		if err != nil {
			latitude = 0
		}
		longitude, err := strconv.ParseFloat(os.Getenv("DEFAULT_LONGITUDE"), 64)
		if err != nil {
			longitude = 0
		}

		// Define a new device model instance.
		newDevice := models.Device{
			DeviceID:    deviceID,
			Name:        "new " + deviceID,
			NetworkType: networkType,
			Latitude:    latitude,
			Longitude:   longitude,
		}

		// Attempt to create a new device record in the database.
		_, err = s.models.Device.Create(&newDevice)
		if err != nil {
			s.errorLog.Printf("Failed to create a new device record for ID %s: %v", deviceID, err)
			continue // Proceed to the next ID instead of stopping.
		}

		// If the network type is 'nb', attempts to create a new device settings record in the database.
		if newDevice.NetworkType == "nb" {
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
