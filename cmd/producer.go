
package main

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/google/uuid"
    amqp "github.com/rabbitmq/amqp091-go"
    "virtflow-go/internal/model"
)

func main() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "connect")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "channel")
    defer ch.Close()

    _, err = ch.QueueDeclare("task.schedule", true, false, false, false, nil)
    failOnError(err, "queue")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    task := model.SchedulingRequest{
        RequestedCPU:    4,
        RequestedMemory: 8192,
        RequestedPool:   "default",
        Dedicated:       false,
        TaskID:          uuid.New().String(),
    }

    body, _ := json.Marshal(task)
    err = ch.PublishWithContext(ctx, "", "task.schedule", false, false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        })
    failOnError(err, "publish")

    log.Printf("Sent task: %s", task.TaskID)
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}
