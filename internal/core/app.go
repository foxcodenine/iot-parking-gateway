package core

import (
	"log"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	AppURL   string
	HttpPort string
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DB       *pgxpool.Pool
	Models   models.Models
	Cache    cache.RedisCache
}
