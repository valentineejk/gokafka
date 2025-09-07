package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Consumer running!")
	})

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"group.id":          "mars",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics([]string{"gods", "^aRegex.*[Tt]opic"}, nil)

	if err != nil {
		panic(err)
	}

	// Run consumer in background goroutine
	go func() {
		defer c.Close()

		for {
			msg, err := c.ReadMessage(time.Second)
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			} else if errKafka, ok := err.(kafka.Error); !ok || errKafka.Code() != kafka.ErrTimedOut {
				fmt.Printf("Consumer error: %v\n", err)
			}
		}
	}()

	log.Fatal(app.Listen(":3000"))

}
