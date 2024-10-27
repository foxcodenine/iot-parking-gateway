package main

import (
	"fmt"
	"log"
	"os"

	"github.com/foxcodenine/iot-parking-gateway/internal/appconfig"
	"github.com/joho/godotenv"
)

var app appconfig.AppConfig

func main() {

	fmt.Println(os.Getenv("TEST"))

	app.Serve()
}

func init() {
	// err := godotenv.Load(".env.development")
	err := godotenv.Load("/app/.env")

	if err != nil {
		fmt.Println(err)
		log.Fatalln("Error loading .env file")
	}
}
