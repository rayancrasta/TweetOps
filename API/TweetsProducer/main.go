package main

import (
	"fmt"
	"log"
	"net/http"
	"tweetsproducer/initializers"
)

const webPort = "8083"

type Config struct {
}

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	app := Config{}

	log.Printf("Starting Users service on port: %s", webPort)

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
