package main

import (
	"fmt"
	"log"
	"net/http"
	"searchquery/initializers"
)

const webPort = "8084"

type Config struct {
}

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	app := Config{}

	log.Printf("Starting Search service on port: %s", webPort)

	// HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	//Start the web server
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}
