package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dy850078/virtflow-scheduler-go/internal/model"

	"github.com/rabbitmq/amqp091-go"
)

func PublishTask(amqpURL string, req model.SchedulingRequest) error {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task.schedule",
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"", q.Name, false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		// Distinguish timeout errors
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("[ERROR] Timeout while publishing task %s", req.TaskID)
			return fmt.Errorf("publish timeout: %w", err)
		}

		log.Printf("[ERROR] Failed to publish task %s: %v", req.TaskID, err)
		return fmt.Errorf("publish: %w", err)
	}

	log.Printf("[DEBUG] Successfully published task: %s", req.TaskID)
	return nil

}
