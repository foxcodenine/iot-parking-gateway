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

var RedisCacheRepo *RedisCache

func NewCache(c *redis.Pool, p string) *RedisCache {
	RedisCacheRepo = &RedisCache{
		Conn:   c,
		Prefix: p,
	}
	return RedisCacheRepo
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

// Set stores a key-value pair in Redis with an optional TTL.
func (rc *RedisCache) Set(key string, value any, ttlSeconds int) error {
	conn := rc.Conn.Get()
	defer conn.Close()

	// Marshal the value to JSON for storage
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	// Set the key with the prefixed key
	if ttlSeconds > 0 {
		// Set with expiration if TTL is specified
		_, err = conn.Do("SETEX", rc.Prefix+key, ttlSeconds, jsonData)
	} else {
		// Set without expiration
		_, err = conn.Do("SET", rc.Prefix+key, jsonData)
	}

	if err != nil {
		return fmt.Errorf("failed to set key in Redis: %w", err)
	}

	return nil
}

// Exists checks if a key exists in Redis.
func (rc *RedisCache) Exists(key string) (bool, error) {
	conn := rc.Conn.Get()
	defer conn.Close()

	// Check if the key exists
	exists, err := redis.Bool(conn.Do("EXISTS", rc.Prefix+key))
	if err != nil {
		return false, fmt.Errorf("failed to check key existence in Redis: %w", err)
	}

	return exists, nil
}

// Get retrieves a key's value from Redis or nil if not found.
func (rc *RedisCache) Get(key string) (any, error) {
	conn := rc.Conn.Get()
	defer conn.Close()

	// Get the value as a string
	data, err := redis.String(conn.Do("GET", rc.Prefix+key))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil // Key does not exist
		}
		return nil, fmt.Errorf("failed to get key from Redis: %w", err)
	}

	// Unmarshal the JSON data into an `any` type
	var value any
	if err := json.Unmarshal([]byte(data), &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return value, nil
}

// Delete removes a key from Redis.
func (rc *RedisCache) Delete(key string) error {
	conn := rc.Conn.Get()
	defer conn.Close()

	// Delete the key
	_, err := conn.Do("DEL", rc.Prefix+key)
	if err != nil {
		return fmt.Errorf("failed to delete key from Redis: %w", err)
	}

	return nil
}
