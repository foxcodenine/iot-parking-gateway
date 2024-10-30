package config

import (
	"log"
)

type App struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
