package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gofiber/fiber/v2"
)

type MessageState int

const (
	MessageCompleted MessageState = iota
	MessageProcessing
	MessageFailed
)

type Message struct {
	State MessageState `json:"state"`
}

func main() {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:29092"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("started kafka")

	defer p.Close()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Producer running!")
	})

	app.Post("/k-producer", func(c *fiber.Ctx) error {

		topic := "gods"
		for i := range 10000 {

			state := MessageState(rand.Intn(3))

			msg := Message{
				State: state,
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Fatal(err)
			}

			//send message
			err = p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          data,
			}, nil)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("sent: %d 🚀", i)
			time.Sleep(1500 * time.Millisecond)

		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status":  true,
			"message": "batch on the way to 🚀 kafka",
		})
	})

	log.Fatal(app.Listen(":4000"))

}
