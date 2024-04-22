package main

import (
	"ConsumerES/initializers"
	"log"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	log.Println("ConsumerES started")

	go consumeMessages()

	//Keep the main thread alive
	select {}
}
