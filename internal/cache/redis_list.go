package cache

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// LPush inserts a new item at the start of the Redis list specified by key.
// The value is serialized to JSON before being inserted.
func (rc *RedisCache) LPush(key string, value any) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}
	conn := rc.Conn.Get()
	defer conn.Close()

	// Prefix the key
	prefixedKey := rc.Prefix + key

	// Push the JSON data to the beginning of the list
	_, err = conn.Do("LPUSH", prefixedKey, jsonData)
	if err != nil {
		return fmt.Errorf("failed to push to Redis list: %v", err)
	}

	return nil
}

// RPush appends a new item to the end of the Redis list specified by key.
// The item is serialized to JSON before being appended.
func (rc *RedisCache) RPush(key string, value any) error {
	// Marshal the value to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Get a Redis connection from the pool
	conn := rc.Conn.Get()
	defer conn.Close()

	// Prefix the key (if a prefix is set) and push the item to the end of the list
	prefixedKey := rc.Prefix + key
	_, err = conn.Do("RPUSH", prefixedKey, jsonData)
	if err != nil {
		return fmt.Errorf("failed to RPUSH to Redis: %w", err)
	}

	return nil
}

// LIndex retrieves an item by index from a Redis list stored under the specified key.
func (rc *RedisCache) LIndex(key string, index int) (any, error) {
	conn := rc.Conn.Get()
	defer conn.Close()

	prefixedKey := rc.Prefix + key

	// Retrieve the item at the specified index as a string
	item, err := redis.String(conn.Do("LINDEX", prefixedKey, index))
	if err != nil {
		if err == redis.ErrNil {
			return nil, fmt.Errorf("item at index %d not found in list: %v", index, err)
		}
		return nil, fmt.Errorf("failed to retrieve item from Redis list: %v", err)
	}

	// Unmarshal the JSON data into an `any` type
	var result any
	err = json.Unmarshal([]byte(item), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from Redis list: %v", err)
	}

	return result, nil
}

// LRange retrieves all items from the Redis list specified by key without deleting the list.
// Each item is deserialized from JSON to an `any` type.
func (rc *RedisCache) LRange(key string) ([]any, error) {
	conn := rc.Conn.Get()
	defer conn.Close()

	prefixedKey := rc.Prefix + key

	// Retrieve all items in the list as strings
	items, err := redis.Strings(conn.Do("LRANGE", prefixedKey, 0, -1))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve items from Redis list: %v", err)
	}

	// Unmarshal each JSON item into an `any` type
	var results []any
	for _, item := range items {
		var value any
		err = json.Unmarshal([]byte(item), &value)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item from Redis list: %v", err)
		}
		results = append(results, value)
	}

	return results, nil
}

// LRangeAndDelete retrieves all items from the Redis list specified by key and then deletes the list.
// Each item is deserialized from JSON to an `any` type.
func (rc *RedisCache) LRangeAndDelete(key string) ([]any, error) {

	rc.mu.Lock()         // Lock before accessing the list
	defer rc.mu.Unlock() // Ensure it's unlocked after

	conn := rc.Conn.Get()
	defer conn.Close()

	prefixedKey := rc.Prefix + key

	// Retrieve all items in the list as strings
	items, err := redis.Strings(conn.Do("LRANGE", prefixedKey, 0, -1))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve items from Redis list: %w", err)
	}

	// Delete the list from Redis
	_, err = conn.Do("DEL", prefixedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to delete Redis list: %w", err)
	}

	// Unmarshal each JSON item into an `any` type
	var results []any
	for _, item := range items {
		var value any
		err = json.Unmarshal([]byte(item), &value)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item from Redis list: %w", err)
		}
		results = append(results, value)
	}

	return results, nil
}
