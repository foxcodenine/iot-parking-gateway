package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// GetDevice retrieves all fields of a device from the Redis hash.
// It returns a map with the device's data or an error if the operation fails.
func (rc *RedisCache) GetDevice(deviceID string) (map[string]any, error) {
	conn := rc.Conn.Get()
	defer conn.Close()

	// Construct the Redis hash key using the device ID.
	hashKey := fmt.Sprintf("%s%s:%s", rc.Prefix, "parking:device", deviceID)

	// Retrieve all fields from the hash using HGETALL.
	data, err := redis.StringMap(conn.Do("HGETALL", hashKey))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve device data for %s: %w", deviceID, err)
	}

	if len(data) == 0 {
		return nil, nil // Return nil if no data exists for the device
	}

	// Deserialize JSON-encoded fields into a map of any type.
	deviceData := make(map[string]any)
	for key, rawValue := range data {
		var parsedValue any
		if err := json.Unmarshal([]byte(rawValue), &parsedValue); err != nil {
			parsedValue = rawValue
		}
		deviceData[key] = parsedValue
	}

	return deviceData, nil
}

// DeleteDevice removes a device from Redis by its ID.
// It returns an error if the deletion fails.
func (rc *RedisCache) DeleteDevice(deviceID string) error {
	conn := rc.Conn.Get()
	defer conn.Close()

	// Construct the Redis hash key using the device ID.
	hashKey := fmt.Sprintf("%s%s:%s", rc.Prefix, "parking:device", deviceID)

	// Delete the hash key from Redis.
	_, err := conn.Do("DEL", hashKey)
	if err != nil {
		return fmt.Errorf("failed to delete device %s: %w", deviceID, err)
	}

	return nil
}

// DeleteAllDevices removes all devices from Redis that match the key pattern "parking:device:*".
// It performs a scan and delete operation and returns an error if any step fails.
func (rc *RedisCache) DeleteAllDevices() error {
	conn := rc.Conn.Get()
	defer conn.Close()

	// Define the key pattern for devices.
	keyPattern := fmt.Sprintf("%s%s:*", rc.Prefix, "parking:device")

	var cursor int64
	for {
		// Use SCAN to find matching keys.
		values, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", keyPattern, "COUNT", 100))
		if err != nil {
			return fmt.Errorf("failed to scan Redis keys: %w", err)
		}

		// Parse SCAN results: cursor and keys.
		cursor, _ = redis.Int64(values[0], nil)
		keys, _ := redis.Strings(values[1], nil)

		// Delete each key found.
		if len(keys) > 0 {
			_, err = conn.Do("DEL", redis.Args{}.AddFlat(keys)...)
			if err != nil {
				return fmt.Errorf("failed to delete device keys: %w", err)
			}
		}

		// If cursor is 0, the scan is complete.
		if cursor == 0 {
			break
		}
	}

	return nil
}

// ProcessParkingEventData updates specific fields for a device in Redis by its ID.
// It handles fields like firmware_version, beacons, happened_at, and is_occupied.
func (rc *RedisCache) ProcessParkingEventData(deviceID string, firmwareVersion string, beacons any, happenedAt string, isOccupied bool) error {
	conn := rc.Conn.Get()
	defer conn.Close()
	// parking:device:

	const dateFormat = "2006-01-02T15:04:05Z" // Define the timestamp format.

	updatedAt := time.Now().UTC().Format(dateFormat) // Format the current time.

	// Construct the Redis hash key.
	hashKey := fmt.Sprintf("%s%s:%s", rc.Prefix, "parking:device", deviceID)

	args := []any{hashKey} // Prepare fields for update.

	// Serialize the `beacons` field to JSON.
	beaconsJSON, err := json.Marshal(beacons)
	if err != nil {
		return fmt.Errorf("failed to marshal beacons field: %v", err)
	}

	// Append key-value pairs for update.
	args = append(args, "firmware_version", firmwareVersion)
	args = append(args, "beacons", beaconsJSON)
	args = append(args, "happened_at", happenedAt)
	args = append(args, "is_occupied", isOccupied)
	args = append(args, "updated_at", updatedAt)

	redisKey := fmt.Sprintf("parking:device:%s", deviceID)

	inCache, _ := rc.Exists(redisKey)
	if !inCache {
		args = append(args, "settings_at", "0001-01-01T00:00:00Z")
		args = append(args, "keepalive_at", "0001-01-01T00:00:00Z")
	}

	// Execute HMSET to update the fields.
	if _, err := conn.Do("HMSET", args...); err != nil {
		return fmt.Errorf("failed to update fields for device %s: %w", deviceID, err)
	}

	return nil
}

func (rc *RedisCache) UpdateKeepaliveAt(deviceID, keepaliveAt, happenedAt, settingsAt string) error {
	conn := rc.Conn.Get()
	defer conn.Close()

	const dateFormat = "2006-01-02T15:04:05Z" // Define the timestamp format.

	updatedAt := time.Now().UTC().Format(dateFormat) // Format the current time.

	// Construct the Redis hash key.
	hashKey := fmt.Sprintf("%s%s:%s", rc.Prefix, "parking:device", deviceID)

	args := []any{hashKey} // Prepare fields for update.

	// Append key-value pairs for update.
	args = append(args, "keepalive_at", keepaliveAt)
	args = append(args, "happened_at", happenedAt)
	args = append(args, "settings_at", settingsAt)
	args = append(args, "updated_at", updatedAt)

	// Execute HMSET to update the fields.
	if _, err := conn.Do("HMSET", args...); err != nil {
		return fmt.Errorf("failed to update fields for device %s: %w", deviceID, err)
	}

	return nil
}

func (rc *RedisCache) UpdateSettingsAt(deviceID, settingsAt, happenedAt, keepaliveAt string) error {
	conn := rc.Conn.Get()
	defer conn.Close()

	const dateFormat = "2006-01-02T15:04:05Z" // Define the timestamp format.

	updatedAt := time.Now().UTC().Format(dateFormat) // Format the current time.

	// Construct the Redis hash key.
	hashKey := fmt.Sprintf("%s%s:%s", rc.Prefix, "parking:device", deviceID)

	args := []any{hashKey} // Prepare fields for update.

	// Append key-value pairs for update.
	args = append(args, "settings_at", settingsAt)
	args = append(args, "happened_at", happenedAt)
	args = append(args, "keepalive_at", keepaliveAt)
	args = append(args, "updated_at", updatedAt)

	// Execute HMSET to update the fields.
	if _, err := conn.Do("HMSET", args...); err != nil {
		return fmt.Errorf("failed to update fields for device %s: %w", deviceID, err)
	}

	return nil
}

// SaveDeviceData saves device data to a Redis hash with a key structured as parking:device:<device_id>.
// It formats date-time fields to ensure consistency across the database and cache.
func (rc *RedisCache) SaveDeviceData(deviceID string, data map[string]any) error {
	const dateFormat = "2006-01-02T15:04:05Z"
	formatDateTimeFields(data, dateFormat) // Ensure date-time fields are formatted.

	conn := rc.Conn.Get()
	defer conn.Close()

	hashKey := fmt.Sprintf("%s%s:%s", rc.Prefix, "parking:device", deviceID)

	for key, value := range data {
		jsonData, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %v", key, err)
		}

		if _, err := conn.Do("HSET", hashKey, key, jsonData); err != nil {
			return fmt.Errorf("failed to HSET data for device %s: %w", deviceID, err)
		}
	}

	return nil
}

// SaveMultipleDevices saves multiple devices to Redis in a single operation.
// Each device's data must include a "device_id" to construct the hash key.
func (rc *RedisCache) SaveMultipleDevices(devices []map[string]any) error {
	const dateFormat = "2006-01-02T15:04:05Z"
	conn := rc.Conn.Get()
	defer conn.Close()

	for _, device := range devices {
		deviceID, ok := device["device_id"].(string)
		if !ok || deviceID == "" {
			return fmt.Errorf("missing or invalid 'device_id' in device: %v", device)
		}

		formatDateTimeFields(device, dateFormat) // Format all date-time fields before saving.
		hashKey := fmt.Sprintf("%s%s:%s", rc.Prefix, "parking:device", deviceID)

		args := []any{hashKey}
		for key, value := range device {
			if key == "device_id" {
				continue
			}

			jsonData, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal value for key %s in device %s: %v", key, deviceID, err)
			}
			args = append(args, key, jsonData)
		}

		if _, err := conn.Do("HMSET", args...); err != nil {
			return fmt.Errorf("failed to save device %s: %w", deviceID, err)
		}
	}

	return nil
}

// UpdateDeviceFields updates multiple key-value pairs in the Redis hash for a device.
// This ensures all timestamp fields are formatted correctly before saving.
func (rc *RedisCache) UpdateDeviceFields(deviceID string, fields map[string]any) error {
	const dateFormat = "2006-01-02T15:04:05Z" // Use a simplified date format

	// Format all date-time fields.
	formatDateTimeFields(fields, dateFormat)

	if _, exists := fields["updated_at"]; !exists {
		fields["updated_at"] = time.Now().UTC().Format(dateFormat)
	}

	conn := rc.Conn.Get()
	defer conn.Close()

	hashKey := fmt.Sprintf("%s%s:%s", rc.Prefix, "parking:device", deviceID)

	args := []any{hashKey}
	for key, value := range fields {
		jsonData, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s in device %s: %v", key, deviceID, err)
		}
		args = append(args, key, jsonData)
	}

	if _, err := conn.Do("HMSET", args...); err != nil {
		return fmt.Errorf("failed to update fields for device %s: %w", deviceID, err)
	}

	return nil
}

// Helper function to format specific date-time fields in the data map.
// This is used to ensure all timestamps are consistent and follow the specified format.
func formatDateTimeFields(data map[string]any, format string) {
	timeFields := []string{"happened_at", "deleted_at", "created_at", "updated_at"}
	for _, field := range timeFields {
		if rawTime, exists := data[field]; exists {
			switch v := rawTime.(type) {
			case string:
				if parsedTime, err := time.Parse(time.RFC3339Nano, v); err == nil {
					data[field] = parsedTime.UTC().Format(format)
				}
			case time.Time:
				data[field] = v.UTC().Format(format)
			}
		}
	}
}
