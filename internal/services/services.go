package services

import (
	"fmt"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
)

type Service struct {
	models *models.Models
	cache  *cache.RedisCache
}

func NewService(m *models.Models, rc *cache.RedisCache) *Service {
	return &Service{
		models: m,
		cache:  rc,
	}
}

// RedisToPostgresRaw retrieves raw data from Redis, saves it to PostgreSQL, and clears the Redis list.
func (s *Service) RedisToPostgresRaw() {
	rawDataLogs, _ := s.cache.LRangeAndDelete("raw-data-logs")

	fmt.Println(rawDataLogs)
}
