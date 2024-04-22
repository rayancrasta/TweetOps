package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

func (app *Config) HandleIncomingTweet(w http.ResponseWriter, r *http.Request) {

	var tweet Tweet
	//Decode the payload
	err := json.NewDecoder(r.Body).Decode(&tweet)
	if err != nil {
		errorJSON(w, fmt.Errorf("error: Failed to parse tweet form: %v", err), http.StatusInternalServerError)
		return
	}

	//Send the request to the producer function
	lang, err := produceMessage(tweet)
	if err != nil {
		errorJSON(w, fmt.Errorf("failed to send reservation request to Kafka: %v", err), http.StatusInternalServerError)
		return
	}

	//Success message
	var response jsonResponse
	response.Message = fmt.Sprintf("Tweet added to %v Topic", lang)
	writeJSON(w, http.StatusOK, response)
}

func produceMessage(tweet Tweet) (string, error) {
	log.Println("DEBUG: Inside Producer_produceMessage ")
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBroker")})

	if err != nil {
		log.Println(err)
		return "", err
	}

	defer producer.Close()

	//Make it byte
	message, err := json.Marshal(tweet)
	if err != nil {
		log.Println(err)
		return "", err
	}

	//Get language
	lang := tweet.Lang

	var topic string

	switch lang {
	case "en":
		topic = os.Getenv("tweetsENTopic")
	case "es":
		topic = os.Getenv("tweetsESTopic")
	default:
		topic = os.Getenv("tweetsRWTopic") //Others
	}

	deliveryChan := make(chan kafka.Event)

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny},
		Value: message,
	}, deliveryChan)

	if err != nil {
		log.Println(err)
		return "", err
	}

	// Wait for delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	// Wait for any outstanding messages to be delivered
	producer.Flush(3 * 1000) // 15-second timeout, adjust as needed

	return lang, nil
}
