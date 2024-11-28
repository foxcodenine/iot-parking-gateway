package core

import (
	"context"
	"errors"
	"log"

	"github.com/foxcodenine/iot-parking-gateway/internal/api/rest/middleware"
	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/foxcodenine/iot-parking-gateway/internal/mq"
	"github.com/foxcodenine/iot-parking-gateway/internal/services"
	"github.com/foxcodenine/iot-parking-gateway/internal/udp"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type App struct {
	AppURL     string
	HttpPort   string
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
	FatalLog   *log.Logger
	DB         *pgxpool.Pool
	Models     models.Models
	MQProducer *mq.RabbitMQProducer
	Cache      *cache.RedisCache
	Cron       *cron.Cron
	UdpServer  *udp.UDPServer
	Service    *services.Service
}

func (app *App) GetUserFromContext(ctx context.Context) (*middleware.UserClaims, error) {
	userData, ok := ctx.Value(middleware.UserContextKey).(*middleware.UserClaims)
	if !ok || userData == nil {
		return nil, errors.New("user claims not found or wrong type in context")
	}
	return userData, nil
}
