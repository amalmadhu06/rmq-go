package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {

	//define rabbitMQ server URL
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// create a new RabbitMQ connection
	rabbitMQConnection, err := amqp.Dial(amqpServerURL)
	if err != nil {
		fmt.Printf("failed to create connection : %v\n", err)
		return
	}

	// close the connection after operations are completed
	defer func(rabbitMQConnection *amqp.Connection) {
		err := rabbitMQConnection.Close()
		if err != nil {
			panic(err)
		}
	}(rabbitMQConnection)

	/* Let's start by opening a channelRMQ to our RabbitMQ
	instance over the connection we have already
	established. */
	channelRMQ, err := rabbitMQConnection.Channel()
	if err != nil {
		fmt.Printf("failed to open a channelRMQ with rabbitMQ connection :%v\n ", err)
		return
	}

	// close the channelRMQ after operations are completed
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {
			panic(err)
		}
	}(channelRMQ)

	// declare Queues
	if _, err := channelRMQ.QueueDeclare(
		"QueueService1", //queue name
		true,            // durable
		false,           //autoDelete
		false,           //exclusive
		false,           //noWait
		nil,             //args
	); err != nil {
		fmt.Printf("failed to declare a queue: %v", err)
		return
	}

	//create a new Fiber instance

	app := fiber.New()

	// add a simple logger middleware
	app.Use(logger.New())

	//routes
	app.Get("/send", func(c *fiber.Ctx) error {
		//create a message to publish

		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(c.Query("msg")),
		}

		//try to publish message to the queue
		if err := channelRMQ.Publish(
			"",
			"QueueService1",
			false,
			false,
			message,
		); err != nil {
			fmt.Printf("failed to publish message : %v", err)
			return nil
		}
		return nil
	})

	// Start Fiber API server.
	log.Fatal(app.Listen(":3000"))
}
