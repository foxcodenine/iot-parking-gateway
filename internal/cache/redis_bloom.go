package cache

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

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
