package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dy850078/virtflow-scheduler-go/internal/algorithm"
	"github.com/dy850078/virtflow-scheduler-go/internal/db"
	"github.com/dy850078/virtflow-scheduler-go/internal/model"
	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("[FATAL] %s: %v", msg, err)
	}
}

func main() {
	connStr := os.Getenv("POSTGRES_CONN")
	if connStr == "" {
		log.Fatal("Missing POSTGRES_CONN env var")
	}

	store, err := db.NewPostgresStore(connStr)
	failOnError(err, "connect to Postgres")
	err = store.InitSchema()
	failOnError(err, "init schema")

	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "open RabbitMQ channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"task.schedule", // name
		true,            // durable
		false,           // auto-delete
		false,           // exclusive
		false,           // no-wait
		nil,             // args
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume("task.schedule", "", true, false, false, false, nil)
	failOnError(err, "consume")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	strategy := &algorithm.CPUStrategy{}
	go consumeLoop(ctx, msgs, strategy, store)

	log.Println("[INFO] Worker running. Ctrl+C to stop.")
	<-ctx.Done()
	log.Println("[INFO] Shutdown complete.")
}

func consumeLoop(ctx context.Context, msgs <-chan amqp091.Delivery, strategy algorithm.SchedulingStrategy, store *db.TaskStore) {
	for {
		select {
		case <-ctx.Done():
			log.Println("[INFO] Context cancelled. Exiting loop.")
			return
		case d := <-msgs:
			var req model.SchedulingRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("[WARN] Failed to parse task: %v", err)
				continue
			}

			nodes := model.MockNodes()
			selected := strategy.SelectBestNode(req, nodes)
			var status, selectedNode string

			if selected != nil {
				status = "success"
				selectedNode = selected.Name
				log.Printf("[INFO] Task %s scheduled to node %s [pool=%s, cpu=%d, mem=%d]",
					req.TaskID, selected.Name, req.RequestedPool, req.RequestedCPU, req.RequestedMemory)
			} else {
				status = "failed"
				log.Printf("[WARN] No suitable node for task %s", req.TaskID)
			}
			if err := store.UpdateStatus(req.TaskID, status, selectedNode); err != nil {
				log.Printf("[ERROR] Failed to update status: %v", err)
			}
		}
	}
}
