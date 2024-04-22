package main

import (
	"ConsumerRW/initializers"
	"log"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	log.Println("ConsumerRW started")

	go consumeMessages()

	//Keep the main thread alive
	select {}
}
