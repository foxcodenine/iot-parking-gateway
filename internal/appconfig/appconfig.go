package appconfig

import (
	"log"
)

type AppConfig struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
