package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/dy850078/virtflow-scheduler-go/internal/model"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func generateTaskID() string {
	return uuid.New().String()
}
func main() {
	log.Println("Start connection...")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"task.schedule", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	request := model.SchedulingRequest{
		RequestedCPU:    4,
		RequestedMemory: 8192,
		RequestedPool:   "default",
		Dedicated:       false,
		TaskID:          generateTaskID(),
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Request cannot convert to JSON")
	}

	body := jsonRequest
	err = ch.PublishWithContext(ctx,
		"",              // exchange
		"task.schedule", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)

}
