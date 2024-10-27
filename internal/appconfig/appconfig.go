package appconfig

import (
	"fmt"
	"log"
	"net/http"
)

type AppConfig struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (app *AppConfig) Serve() {

	webPort := "8080"
	fmt.Println("Starting web server - http://localhost:" + webPort + "/")

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}
}
