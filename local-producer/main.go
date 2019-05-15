package main

import (
	"fmt"
	"math/rand"
	"time"
	"encoding/json"
	"log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type OrderStatusMessage struct {
	EventType             string `json:"eventType"`
	OrderID               int    `json:"orderId"`
	RestaurantName        string `json:"restaurantName"`
	NotificationID        string `json:"notificationId"`
	Refr                  string `json:"refr"`
	OrderType             string `json:"orderType"`
	Language              string `json:"language"`
	EstimatedTimeDelivery string `json:"estimatedTimeDelivery"`
	RequestedDeliveryTime string `json:"requestedDeliveryTime"`
}

func main() {
	println("local-producer started...")

	//Create our correct sequence
	correctSequence := createMessageCorrectSequence()

	channel := make(chan bool)

	//Starts
	go sendCorrectSequence(correctSequence, channel)

	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		fmt.Println("Ticker tick...")

		//Create our correct sequence
		correctSequence := createMessageCorrectSequence()
		//Create our incorrect sequence
		incorrectSequence := createMessageIncorrectSequence()

		//This triggers sendCorrectSequence because of the lock <-channel
		channel <- false
		if <-channel {
			fmt.Println("true")
			go sendIncorrectSequence(incorrectSequence, channel)
		} else {
			fmt.Println("false")
			go sendCorrectSequence(correctSequence, channel)
		}
	}
}

func sendCorrectSequence(correctSequence []OrderStatusMessage, channel chan bool) {
	<-channel

	fmt.Println("CorrectChannel")
	fmt.Println("Send Correct Sequence starting...")

		p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9093"})
		if err != nil {
			log.Fatal(err)
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

		for _, message := range correctSequence {
			topic := "orderStatus"
			data, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
				return
			}

			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte(data),
				Key: 			[]byte("QWERTY_1"),
			}, nil)
		}

		// Wait for message deliveries before shutting down
		p.Flush(15 * 1000)
		println("ending producer...")

	//TOGGLE
	channel <- true
}

func sendIncorrectSequence(incorrectSequence []OrderStatusMessage, channel chan bool) {
	<-channel
	fmt.Println("IncorrectChannel")
	fmt.Println("Send Incorrect Sequence starting...")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9093"})
	if err != nil {
		log.Fatal(err)
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

	for _, message := range incorrectSequence {
		topic := "orderStatus"
		data, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(data),
			Key: 			[]byte("QWERTY_2"),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
	println("ending producer...")

	//TOGGLE
	channel <- false
}


func createMessageCorrectSequence() []OrderStatusMessage {

	rand.Seed(time.Now().UnixNano())
	deterministicRandomOrderID := rand.Intn(100000)
	var correctSequence []OrderStatusMessage

	//Create Correct Sequence
	for i := 0; i <= 5; i++ {
		var message OrderStatusMessage
		switch i {
		case 0:
			message.EventType = "NEW"
			message.OrderID = deterministicRandomOrderID
			message.RestaurantName = "Test Restaurant"
			message.NotificationID = "12345-34812-23333:434343-233"
			message.Refr = "com.yourdelivery.lieferando"
			message.OrderType = "DELIVERY"
			message.Language = "en"
			message.RequestedDeliveryTime = "text_asap"
			correctSequence = append(correctSequence, message)
			// fmt.Println("one")
		case 1:
			message.EventType = "CONFIRMED"
			message.RestaurantName = "Test Restaurant"
			message.OrderID = deterministicRandomOrderID
			message.EstimatedTimeDelivery = "385"
			message.OrderType = "DELIVERY"
			message.RequestedDeliveryTime = "text_asap"
			correctSequence = append(correctSequence, message)
			// fmt.Println("two")
		case 2:
			message.EventType = "KITCHEN"
			message.RestaurantName = "Test Restaurant"
			message.OrderID = deterministicRandomOrderID
			message.OrderType = "DELIVERY"
			message.RequestedDeliveryTime = "text"
			correctSequence = append(correctSequence, message)
			// fmt.Println("three")
		case 4:
			message.EventType = "ENROUTE"
			message.RestaurantName = "Test Restaurant"
			message.OrderID = deterministicRandomOrderID
			message.OrderType = "DELIVERY"
			message.RequestedDeliveryTime = "text"
			correctSequence = append(correctSequence, message)
			// fmt.Println("four")
		case 5:
			message.EventType = "TIMECHANGE"
			message.RestaurantName = "Test Restaurant"
			message.OrderID = deterministicRandomOrderID
			message.EstimatedTimeDelivery = "35"
			message.OrderType = "DELIVERY"
			message.RequestedDeliveryTime = "text_asap"
			correctSequence = append(correctSequence, message)
			// fmt.Println("five")
		}
	}

	return correctSequence
}


func createMessageIncorrectSequence() []OrderStatusMessage {
	rand.Seed(time.Now().UnixNano())
	deterministicRandomOrderID := rand.Intn(10000)
	var correctSequence []OrderStatusMessage

	//Create Incorrect Sequence
	for i := 0; i <= 5; i++ {
		var message OrderStatusMessage
		switch i {
		case 0:
			message.EventType = "CONFIRMED"
			message.OrderID = deterministicRandomOrderID
			message.EstimatedTimeDelivery = "358"
			message.OrderType = "DELIVERY"
			message.RequestedDeliveryTime = "text_asap"
			correctSequence = append(correctSequence, message)
			// fmt.Println("two")
		case 1:
			message.EventType = "NEW"
			message.OrderID = deterministicRandomOrderID
			message.RestaurantName = "Test Restaurant"
			message.NotificationID = "12345-34812-23333:434343-233"
			message.Refr = "com.yourdelivery.lieferando"
			message.OrderType = "DELIVERY"
			message.Language = "en"
			message.RequestedDeliveryTime = "text"
			correctSequence = append(correctSequence, message)
			// fmt.Println("one")
		case 2:
			message.EventType = "KITCHEN"
			message.OrderID = deterministicRandomOrderID
			message.OrderType = "DELIVERY"
			message.RequestedDeliveryTime = "text"
			correctSequence = append(correctSequence, message)
			// fmt.Println("three")
		case 4:
			message.EventType = "ENROUTE"
			message.OrderID = deterministicRandomOrderID
			message.OrderType = "DELIVERY"
			message.RequestedDeliveryTime = "text"
			correctSequence = append(correctSequence, message)
			// fmt.Println("four")
		case 5:
			message.EventType = "TIMECHANGE"
			message.OrderID = deterministicRandomOrderID
			message.EstimatedTimeDelivery = "35"
			message.OrderType = "DELIVERY"
			message.RequestedDeliveryTime = "text_timed"
			correctSequence = append(correctSequence, message)
			// fmt.Println("five")
		}
	}

	return correctSequence
}