package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type SafeTweetCheck struct {
	Status_id int    `json:"status_id"`
	Text      string `json:"text"`
}

func safeCheck(tweet Tweet) error {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBroker")})

	if err != nil {
		log.Println(err)
		return err
	}

	defer producer.Close()

	//Body that is sent to safechecking
	var safetweetcheck SafeTweetCheck
	safetweetcheck.Status_id = tweet.Status_id
	safetweetcheck.Text = tweet.Text

	//Make it byte
	message, err := json.Marshal(safetweetcheck)
	if err != nil {
		log.Println(err)
		return err
	}

	topic := "safecheck"

	deliveryChan := make(chan kafka.Event)

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny},
		Value: message,
	}, deliveryChan)

	if err != nil {
		log.Println(err)
		return err
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

	return nil
}
