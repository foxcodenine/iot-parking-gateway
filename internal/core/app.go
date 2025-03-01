package core

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/foxcodenine/iot-parking-gateway/internal/apptypes"
	"github.com/foxcodenine/iot-parking-gateway/internal/cache"
	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/foxcodenine/iot-parking-gateway/internal/mq"
	"github.com/foxcodenine/iot-parking-gateway/internal/services"
	"github.com/foxcodenine/iot-parking-gateway/internal/udp"
	socketio "github.com/googollee/go-socket.io"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type App struct {
	AppURL     string
	HttpPort   string
	DB         *pgxpool.Pool
	Models     models.Models
	MQProducer *mq.RabbitMQProducer
	Cache      *cache.RedisCache
	Cron       *cron.Cron
	UdpServer  *udp.UDPServer
	SocketIO   *socketio.Server

	Service          *services.Service
	DeviceAccessMode *string
}

func (app *App) GetUserFromContext(ctx context.Context) (*apptypes.UserClaims, error) {
	userData, ok := ctx.Value(apptypes.UserContextKey).(*apptypes.UserClaims)
	if !ok || userData == nil {
		return nil, errors.New("user claims not found or wrong type in context")
	}
	return userData, nil
}

func (app *App) PushAuditToCache(
	userData apptypes.UserClaims, // User performing the action
	action string, // Action type (e.g., "UPDATE")
	entity string, // Entity being acted upon
	entityID string, // ID of the entity
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
	app.Cache.RPush("logs:audit-logs", auditLogEntry)
}

// Helper function to get client IP
func getClientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback to returning the whole field
	}
	return ip
}
