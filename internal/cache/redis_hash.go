package cache

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// HSet adds a key-value pair to a Redis hash map. The value is serialized to JSON before being stored.
func (rc *RedisCache) HSet(mapKey string, fieldKey string, value any) error {
	conn := rc.Conn.Get()
	defer conn.Close()

	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	prefixedMapKey := rc.Prefix + mapKey
	_, err = conn.Do("HSET", prefixedMapKey, fieldKey, jsonData)
	if err != nil {
		return fmt.Errorf("failed to HSET in Redis: %w", err)
	}

	return nil
}

// HGet retrieves a value by key from a Redis hash map and deserializes it from JSON to an `any` type.
func (rc *RedisCache) HGet(mapKey string, fieldKey string) (any, error) {
	conn := rc.Conn.Get()
	defer conn.Close()

	prefixedMapKey := rc.Prefix + mapKey
	item, err := redis.String(conn.Do("HGET", prefixedMapKey, fieldKey))
	if err != nil {
		if err == redis.ErrNil {
			return nil, fmt.Errorf("field %s not found in map", fieldKey)
		}
		return nil, fmt.Errorf("failed to retrieve item from Redis hash map: %v", err)
	}

	var result any
	err = json.Unmarshal([]byte(item), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal item from Redis hash map: %v", err)
	}

	return result, nil
}

// HGetAll retrieves all key-value pairs from a Redis hash map and deserializes them from JSON to `any` types.
func (rc *RedisCache) HGetAll(mapKey string) (map[string]any, error) {
	conn := rc.Conn.Get()
	defer conn.Close()

	prefixedMapKey := rc.Prefix + mapKey
	items, err := redis.StringMap(conn.Do("HGETALL", prefixedMapKey))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve items from Redis hash map: %v", err)
	}

	results := make(map[string]any)
	for key, item := range items {
		var value any
		err = json.Unmarshal([]byte(item), &value)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item from Redis hash map: %v", err)
		}
		results[key] = value
	}

	return results, nil
}
