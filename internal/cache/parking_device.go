package cache

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// SaveDeviceData saves device data to a Redis hash with a key structured as parking:device:<device_id>.
func (rc *RedisCache) SaveDeviceData(deviceID string, data map[string]any) error {

	delete(data, "beacons_json")
	data["delete_at"] = nil

	conn := rc.Conn.Get()
	defer conn.Close()

	// Construct the Redis hash key using the device ID.
	hashKey := fmt.Sprintf("%s%s:%s", rc.Prefix, "parking:device", deviceID)

	for key, value := range data {
		jsonData, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %v", key, err)
		}

		// Store each key-value pair in the Redis hash.
		if _, err := conn.Do("HSET", hashKey, key, jsonData); err != nil {
			return fmt.Errorf("failed to HSET data for device %s: %w", deviceID, err)
		}
	}

	return nil
}

// UpdateDeviceFields updates multiple key-value pairs in the Redis hash for a device.
func (rc *RedisCache) UpdateDeviceFields(deviceID string, fields map[string]any) error {

	delete(fields, "beacons_json")
	fields["delete_at"] = nil

	conn := rc.Conn.Get()
	defer conn.Close()

	// Construct the Redis hash key using the device ID.
	hashKey := fmt.Sprintf("%s%s:%s", rc.Prefix, "parking:device", deviceID)

	// Create a slice to hold key-value pairs for the HMSET command.
	args := []any{hashKey}

	// Serialize each value to JSON and append key-value pairs to the arguments slice.
	for key, value := range fields {
		jsonData, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %v", key, err)
		}
		args = append(args, key, jsonData)
	}

	// Update the hash in Redis using HMSET.
	_, err := conn.Do("HMSET", args...)
	if err != nil {
		return fmt.Errorf("failed to update fields for device %s: %w", deviceID, err)
	}

	return nil
}

// GetDevice retrieves all fields of a device from the Redis hash.
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

	// Deserialize JSON strings back to their original types.
	deviceData := make(map[string]any)
	for key, jsonValue := range data {
		var value any
		if err := json.Unmarshal([]byte(jsonValue), &value); err != nil {
			return nil, fmt.Errorf("failed to unmarshal value for key %s: %v", key, err)
		}
		deviceData[key] = value
	}

	return deviceData, nil
}

// DeleteDevice removes a device from Redis by its ID.
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
func (rc *RedisCache) DeleteAllDevices() error {
	conn := rc.Conn.Get()
	defer conn.Close()

	// Define the key pattern for devices
	keyPattern := fmt.Sprintf("%s%s:*", rc.Prefix, "parking:device")

	var cursor int64
	for {
		// Use SCAN to find matching keys
		values, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", keyPattern, "COUNT", 100))
		if err != nil {
			return fmt.Errorf("failed to scan Redis keys: %w", err)
		}

		// Parse SCAN results: cursor and keys
		cursor, _ = redis.Int64(values[0], nil)
		keys, _ := redis.Strings(values[1], nil)

		// Delete each key found
		if len(keys) > 0 {
			_, err = conn.Do("DEL", redis.Args{}.AddFlat(keys)...)
			if err != nil {
				return fmt.Errorf("failed to delete device keys: %w", err)
			}
		}

		// If cursor is 0, the scan is complete
		if cursor == 0 {
			break
		}
	}

	return nil
}

// SaveMultipleDevices saves multiple devices to Redis in a single operation.
func (rc *RedisCache) SaveMultipleDevices(devices map[string]map[string]any) error {
	conn := rc.Conn.Get()
	defer conn.Close()

	for deviceID, fields := range devices {
		delete(fields, "beacons_json")
		fields["delete_at"] = nil

		// Construct the Redis hash key
		hashKey := fmt.Sprintf("%s%s:%s", rc.Prefix, "parking:device", deviceID)

		// Prepare the fields for HMSET
		args := []any{hashKey}
		for key, value := range fields {
			jsonData, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal value for key %s: %v", key, err)
			}
			args = append(args, key, jsonData)
		}

		// Save the hash
		_, err := conn.Do("HMSET", args...)
		if err != nil {
			return fmt.Errorf("failed to save device %s: %w", deviceID, err)
		}
	}

	return nil
}

// ProcessParkingEventData updates specific fields (firmware_version, beacons, happened_at, is_occupied)
// for a device in Redis by its ID.
func (rc *RedisCache) ProcessParkingEventData(deviceID string, firmwareVersion string, beacons any, happenedAt string, isOccupied bool) error {
	conn := rc.Conn.Get()
	defer conn.Close()

	// Construct the Redis hash key
	hashKey := fmt.Sprintf("%s%s:%s", rc.Prefix, "parking:device", deviceID)

	// Prepare the fields to update
	args := []any{hashKey}

	// Serialize the `beacons` field to JSON
	beaconsJSON, err := json.Marshal(beacons)
	if err != nil {
		return fmt.Errorf("failed to marshal beacons field: %v", err)
	}

	// Add key-value pairs for the fields to be updated
	args = append(args, "firmware_version", firmwareVersion)
	args = append(args, "beacons", beaconsJSON)
	args = append(args, "happened_at", happenedAt)
	args = append(args, "is_occupied", isOccupied)

	// Execute HMSET to update the fields
	_, err = conn.Do("HMSET", args...)
	if err != nil {
		return fmt.Errorf("failed to update fields for device %s: %w", deviceID, err)
	}

	return nil
}
