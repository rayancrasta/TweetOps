package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
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

	var KafkaConsumerConfig = &kafka.ConfigMap{
		"bootstrap.servers":  os.Getenv("kafkaBroker"),
		"group.id":           os.Getenv("consumergroupid"),
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "false",
	}

	consumer, err := kafka.NewConsumer(KafkaConsumerConfig)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	// Set up Elasticsearch client
	var ESConfig = elasticsearch.Config{
		Addresses: []string{os.Getenv("elasticsearchaddr")}, // Change this to your Elasticsearch address
	}

	esClient, err := elasticsearch.NewClient(ESConfig)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
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
				req := esapi.IndexRequest{
					Index:      "tweetscombined",
					DocumentID: strconv.Itoa(tweet.Status_id),
					Body:       bytes.NewReader(msg.Value),
					Refresh:    "true",
				}

				res, err := req.Do(context.Background(), esClient)
				if err != nil {
					log.Printf("Error indexing tweet into Elasticsearch: %v", err)
					continue
				}
				defer res.Body.Close()

				//Send for safecheck
				err = safeCheck(tweet)
				if err != nil {
					log.Printf("Error sending tweet into Elasticsearch: %v", err)
					continue
				}

				// Commit the Kafka message offset
				if _, err := consumer.CommitMessage(msg); err != nil {
					log.Printf("Error committing offset: %v", err)
				}

			} else {
				fmt.Printf("Consumer error: %v\n", err)
			}
		}

	}

}
