
package main

import (
    "context"
    "encoding/json"
    "log"
    "virtflow-go/internal/algorithm"
    "virtflow-go/internal/db"
    "virtflow-go/internal/model"

    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "connect")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "channel")
    defer ch.Close()

    msgs, err := ch.Consume("task.schedule", "", true, false, false, false, nil)
    failOnError(err, "consume")

    db.InitDB()

    log.Println(" [*] Waiting for messages. To exit press CTRL+C")
    forever := make(chan bool)

    go func() {
        for d := range msgs {
            var req model.SchedulingRequest
            err := json.Unmarshal(d.Body, &req)
            if err != nil {
                log.Printf("❌ Failed to parse task: %s", err)
                continue
            }

            log.Printf("📦 Task received: %+v", req)
            db.UpdateTaskStatus(req.TaskID, "running", "")

            selected := algorithm.SelectBestNode(req, model.MockNodes())
            if selected != nil {
                db.UpdateTaskStatus(req.TaskID, "success", selected.Name)
                log.Printf("✅ Task %s scheduled to: %s", req.TaskID, selected.Name)
            } else {
                db.UpdateTaskStatus(req.TaskID, "failed", "no available node")
                log.Printf("❌ No node selected for task %s", req.TaskID)
            }
        }
    }()

    <-forever
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}
