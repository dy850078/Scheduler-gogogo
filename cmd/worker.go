package main

import (
	"encoding/json"
	"log"

	"github.com/dy850078/virtflow-scheduler-go/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {

	// while true:
	// get message from RabbitMQ
	// decode JSON payload
	// log task content

	log.Println("Start connection...")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	msgs, err := ch.Consume(
		"task.schedule", // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			var req model.SchedulingRequest
			err := json.Unmarshal(d.Body, &req)
			if err != nil {
				log.Printf(" [Warn] Failed to parse task %s", err)
				continue
			}
			log.Printf(" [Queue] Task Received: %+v", req)
		}
	}()

	log.Printf(" [Queue] Waiting for messages. To exit press CTRL+C")
	select {}

}
