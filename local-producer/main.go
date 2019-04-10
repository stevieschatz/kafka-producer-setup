package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	println("starting producer...")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9093"})
	if err != nil {
		log.Fatal(err)
	}

	var orderPayload []string

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "order-status") {
			orderPayload = append(orderPayload, file.Name())
		}
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	var messages []string

	for _, fileName := range orderPayload {
		jsonFile, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer jsonFile.Close()
		byteValue, err := ioutil.ReadAll(jsonFile)

		if err != nil && err == io.EOF {
			log.Fatal(err)
		}

		var result map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)

		messages = append(messages, string(byteValue))
		fmt.Println(result)

	}

	for range time.NewTicker(10 * time.Second).C {
		for _, message := range messages {
			topic := "orderStatus"
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte(message),
			}, nil)
		}
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
	println("ending producer...")

}
