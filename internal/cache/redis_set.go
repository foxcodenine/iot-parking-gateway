package cache

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func (rc *RedisCache) SAdd(key string, values ...any) error {
	conn := rc.Conn.Get()
	defer conn.Close()

	args := make([]any, len(values)+1)
	args[0] = rc.Prefix + key

	for i, v := range values {
		jsonData, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %v", err)
		}
		args[i+1] = jsonData
	}

	_, err := conn.Do("SADD", args...)
	if err != nil {
		return fmt.Errorf("failed to SADD to Redis set: %w", err)
	}

	return nil
}

func (rc *RedisCache) SMembers(key string) ([]any, error) {
	conn := rc.Conn.Get()
	defer conn.Close()

	memberStrings, err := redis.Strings(conn.Do("SMEMBERS", rc.Prefix+key))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve SMEMBERS from Redis set: %w", err)
	}

	members := make([]any, len(memberStrings))
	for i, memberStr := range memberStrings {
		var member any
		if err := json.Unmarshal([]byte(memberStr), &member); err != nil {
			return nil, fmt.Errorf("failed to unmarshal member: %w", err)
		}
		members[i] = member
	}

	return members, nil
}

func (rc *RedisCache) SMembersDel(key string) ([]any, error) {
	rc.mu.Lock()         // Lock before accessing the set
	defer rc.mu.Unlock() // Ensure it's unlocked after operation

	conn := rc.Conn.Get()
	defer conn.Close()

	// Get the prefixed key
	prefixedKey := rc.Prefix + key

	// Retrieve all items in the set as strings
	memberStrings, err := redis.Strings(conn.Do("SMEMBERS", prefixedKey))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve SMEMBERS from Redis set: %w", err)
	}

	// Delete the set from Redis
	_, err = conn.Do("DEL", prefixedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to delete Redis set after retrieving members: %w", err)
	}

	// Convert the string data into the desired 'any' type using json.Unmarshal
	members := make([]any, len(memberStrings))
	for i, memberStr := range memberStrings {
		var member any
		if err := json.Unmarshal([]byte(memberStr), &member); err != nil {
			return nil, fmt.Errorf("failed to unmarshal member: %w", err)
		}
		members[i] = member
	}

	return members, nil
}
