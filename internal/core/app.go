package core

import (
	"log"

	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/foxcodenine/iot-parking-gateway/internal/services"
	"github.com/foxcodenine/iot-parking-gateway/internal/udp"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type App struct {
	AppURL    string
	HttpPort  string
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
	DB        *pgxpool.Pool
	Models    models.Models
	Cache     *cache.RedisCache
	Cron      *cron.Cron
	UdpServer *udp.UDPServer
	Service   *services.Service
}
