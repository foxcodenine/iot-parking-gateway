package config

import (
	"log"

	"github.com/foxcodenine/iot-parking-gateway/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DB       *pgxpool.Pool
	Models   models.Models
}
