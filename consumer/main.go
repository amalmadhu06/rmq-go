package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {

	// get rabbitMQ server url from env
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	//create a new rabbitMQ connection
	rMQConnection, err := amqp.Dial(amqpServerURL)
	if err != nil {
		fmt.Printf("failed to create rabbitMQ connection : %v\n", err)
		return
	}

	//close the connection after needed
	defer func(rMQConnection *amqp.Connection) {
		if err := rMQConnection.Close(); err != nil {
			fmt.Printf("faied to close rabbitMQ connection : %v\n", err)
			return
		}
	}(rMQConnection)

	// Opening a channel to our RabbitMQ instance over
	// the connection we have already established.
	channelRMQ, err := rMQConnection.Channel()
	if err != nil {
		fmt.Printf("failed to create channel with rabbitMQ connection : %v\n", err)
	}

	// close the channel after operation
	defer func(channelRMQ *amqp.Channel) {
		if err := channelRMQ.Close(); err != nil {
			fmt.Printf("faied to close rabbitMQ channel : %v\n", err)
			return
		}
	}(channelRMQ)

	// subscribe to our QueueService1 for getting the messages
	messages, err := channelRMQ.Consume(
		"QueueService1",
		"",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		fmt.Printf("failed to subscribe to rabbitMQ channel : %v", err)
	}

	// Build a welcome message.
	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	//make a channel to receive messages ino an infinite loop
	forever := make(chan bool)

	go func() {
		for message := range messages {
			// printing the received message for example
			log.Printf("> Received Message : %s\n", message.Body)
		}
	}()
	<-forever
}
