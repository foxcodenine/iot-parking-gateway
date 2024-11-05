package cache

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisCache struct {
	Conn   *redis.Pool
	Prefix string
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
