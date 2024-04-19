package main

import (
	"fmt"
	"log"
	"net/http"
	initializers "relations/initializers"
)

const webPort = "8082"

type Config struct {
}

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	app := Config{}

	log.Printf("Starting Relations service on port: %s", webPort)

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
