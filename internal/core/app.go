package core

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

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

func (app *App) PushAuditToCache(
	userData middleware.UserClaims, // User performing the action
	action string, // Action type (e.g., "UPDATE")
	entity string, // Entity being acted upon
	entityID int, // ID of the entity
	r *http.Request, // HTTP request for context
	details string, // Detailed description of the action

) {
	// Create the audit log entry
	auditLogEntry := models.AuditLog{
		UserID:      userData.UserID,
		Email:       userData.Email,
		AccessLevel: userData.AccessLevel,
		HappenedAt:  time.Now().UTC(),
		Action:      action,
		Entity:      entity,
		EntityID:    entityID,
		URL:         r.URL.Path,
		IPAddress:   getClientIP(r),
		Details:     details,
	}

	// Push the audit log entry to the cache
	app.Cache.RPush("audit-logs", auditLogEntry)
}

// Helper function to get client IP
func getClientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback to returning the whole field
	}
	return ip
}
