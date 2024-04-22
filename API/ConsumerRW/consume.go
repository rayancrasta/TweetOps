package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Tweet struct {
	Status_id          int       `json:"status_id"`
	User_id            int       `json:"user_id"`
	Created_at         time.Time `json:"created_at"`
	Screen_name        string    `json:"screen_name"`
	Text               string    `json:"text"`
	Reply_to_status_id int       `json:"reply_to_status_id"`
	Reply_to_user_id   int       `json:"reply_to_user_id"`
	Favourites_count   int       `json:"favourites_count"`
	Retweet_count      int       `json:"retweet_count"`
	Country_code       string    `json:"country_code"`
	Place_full_name    string    `json:"place_full_name"`
	Place_type         string    `json:"place_type"`
	Verified           bool      `json:"verified"`
	Lang               string    `json:"lang"`
}

func consumeMessages() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  os.Getenv("kafkaBroker"),
		"group.id":           "tweetRW-consumer-group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "true",
	})

	//commit it later on putting tweet in elasticsearch

	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	//Subscribe to Kafka topic
	consumer.SubscribeTopics([]string{os.Getenv("tweetsTopic")}, nil) //checkDoc

	// Signal channel to handle graceful termination in case of interrupts (SIGINT (controlC) or SIGTERM).
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	//Infinite loop to continuosuly consume messages
	for {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal: %v: terminating ", sig)
			consumer.Close()
			return
		default:
			//Read the message from Kafka topic
			msg, err := consumer.ReadMessage(-1) // Timeout parameter. -1 is to wait indenitely for the message
			//block until a new message is available or an error occurs.

			if err == nil { //success
				var tweet Tweet
				err := json.Unmarshal(msg.Value, &tweet)

				if err != nil {
					log.Printf("Error decoding JSON: %v", err)
					continue
				}

				log.Println(tweet)

				//Put into elasticsearch from here

			} else {
				fmt.Printf("Consumer error: %v\n", err)
			}
		}

	}

}
