package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisCache struct {
	Conn   *redis.Pool
	Prefix string
	mu     sync.Mutex
}

// CreateRedisPool initializes and returns a Redis connection pool.
func CreateRedisPool() (*redis.Pool, error) {

	var redisPort string

	switch os.Getenv("GO_ENV") {
	case "production":
		redisPort = os.Getenv("REDIS_PORT")
	default:
		redisPort = os.Getenv("REDIS_PORT_EX")
	}

	// Parse the Redis database index from an environment variable
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0 // Default to DB 0 if parsing fails
	}

	pool := &redis.Pool{
		MaxIdle:     50,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			redisAddress := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), redisPort)
			conn, err := redis.Dial(
				"tcp",
				redisAddress,
				redis.DialPassword(os.Getenv("REDIS_PASSWORD")),
				redis.DialDatabase(redisDB),
			)
			if err != nil {
				return nil, err
			}

			// Test the connection
			if _, err := conn.Do("PING"); err != nil {
				conn.Close()
				return nil, err
			}

			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, lastUsed time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// Test the pool connection by getting a connection and PINGing
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return pool, nil
}

func (rc *RedisCache) LPush(key string, value any) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}
	conn := rc.Conn.Get()
	defer conn.Close()

	// Prefix the key
	prefixedKey := fmt.Sprintf("%s:%s", rc.Prefix, key)

	// Push the JSON data to the beginning of the list
	_, err = conn.Do("LPUSH", prefixedKey, jsonData)
	if err != nil {
		return fmt.Errorf("failed to push to Redis list: %v", err)
	}

	return nil
}

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

// CreateBloomFilter creates a new Bloom Filter in Redis with the specified error rate and capacity.
func (rc *RedisCache) CreateBloomFilter(filterName string, errorRate float64, capacity int) error {
	conn := rc.Conn.Get()
	defer conn.Close()

	_, err := conn.Do("BF.RESERVE", rc.Prefix+filterName, errorRate, capacity)
	if err != nil {
		return fmt.Errorf("failed to create Bloom filter: %v", err)
	}

	return nil
}

// CheckItemInBloomFilter checks if an item exists in the specified Bloom Filter.
func (rc *RedisCache) CheckItemInBloomFilter(filterName string, item string) (bool, error) {
	conn := rc.Conn.Get()
	defer conn.Close()

	// Execute the BF.EXISTS command to check the existence of the item in the Bloom Filter
	exists, err := redis.Int(conn.Do("BF.EXISTS", rc.Prefix+filterName, item))
	if err != nil {
		return false, fmt.Errorf("failed to check item in Bloom filter: %v", err)
	}

	// The BF.EXISTS command returns 1 if the item is possibly in the set, 0 if it is definitely not in the set
	return exists == 1, nil
}

// AddItemToBloomFilter adds an item to the specified Bloom Filter.
func (rc *RedisCache) AddItemToBloomFilter(filterName string, item string) (bool, error) {
	conn := rc.Conn.Get()
	defer conn.Close()

	// Execute the BF.ADD command to add the item to the Bloom Filter
	added, err := redis.Int(conn.Do("BF.ADD", rc.Prefix+filterName, item))
	if err != nil {
		return false, fmt.Errorf("failed to add item to Bloom filter: %v", err)
	}

	// The BF.ADD command returns 1 if the item is a new addition to the filter, 0 if it was already present
	return added == 1, nil
}

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
