package main

import (
	"ConsumerEN/initializers"
	"log"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	log.Println("ConsumerEN started")

	go consumeMessages()

	//Keep the main thread alive
	select {}
}
