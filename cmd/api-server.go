package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/dy850078/virtflow-scheduler-go/internal/model"
	"github.com/dy850078/virtflow-scheduler-go/internal/queue"
	"github.com/dy850078/virtflow-scheduler-go/internal/util"
)

func main() {
	http.HandleFunc("/schedule", handleSchedule)

	log.Println("[INFO] API Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("[FATAL] failed to start server: %v", err)
	}
}

func handleSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.SchedulingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	req.TaskID = util.GenerateTaskID()
	if req.TaskID == "" {
		http.Error(w, "Missing task_id", http.StatusBadRequest)
		return
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	err := queue.PublishTask(rabbitURL, req)
	if err != nil {
		log.Printf("[ERROR] Failed to publish task: %v", err)
		http.Error(w, "Failed to enqueue task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Task accepted"))
}
