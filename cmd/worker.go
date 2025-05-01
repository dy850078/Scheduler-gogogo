package main

import (
	"context"
	"encoding/json"
	"log"
	"os/signal"
	"syscall"

	"github.com/dy850078/virtflow-scheduler-go/internal/algorithm"
	"github.com/dy850078/virtflow-scheduler-go/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

// failOnError handles error with log and panic for fatal errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("[FATAL] %s: %v", msg, err)
	}
}

func main() {
	log.Println("[INFO] Starting RabbitMQ worker...")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume(
		"task.schedule", // queue name
		"",              // consumer tag
		true,            // auto-ack (enable manual ack in future)
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	failOnError(err, "Failed to register a consumer")

	// Create a cancellable context that listens for Ctrl+C or termination
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	strategy := &algorithm.CPUStrategy{} // May change to different strategy

	// Start the consumer loop in a separate goroutine
	go consumeLoop(ctx, msgs, strategy)

	log.Println("[INFO] Worker is running. Press Ctrl+C to stop.")
	<-ctx.Done()
	log.Println("[INFO] Shutdown signal received. Cleaning up...")
}

// consumeLoop handles incoming RabbitMQ messages in a separate goroutine
func consumeLoop(ctx context.Context, msgs <-chan amqp.Delivery, strategy algorithm.SchedulingStrategy) {
	for {
		select {
		case <-ctx.Done():
			log.Println("[INFO] Consumer loop exiting due to context cancellation")
			return

		case d, ok := <-msgs:
			if !ok {
				log.Println("[WARN] RabbitMQ message channel closed")
				return
			}

			var req model.SchedulingRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("[WARN] Failed to parse task payload: %v", err)
				continue
			}

			log.Printf("[INFO] Received task: TaskID=%s CPU=%d Mem=%d Pool=%s Dedicated=%t",
				req.TaskID, req.RequestedCPU, req.RequestedMemory, req.RequestedPool, req.Dedicated)

			nodes := model.MockNodes() // TODO: API/DB Access
			selected := strategy.SelectBestNode(req, nodes)

			if selected != nil {
				log.Printf("[INFO] Task %s scheduled to node: %s", req.TaskID, selected.Name)
			} else {
				log.Printf("[WARN] No suitable node found for task: %s", req.TaskID)
			}
		}
	}
}
