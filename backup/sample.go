package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// BareMetalNode represents a bare metal server
type BareMetalNode struct {
	Name       string  `json:"name"`
	CPU        int     `json:"cpu"`
	Memory     int     `json:"memory"`
	Storage    int     `json:"storage"`
	UsageCPU   float64 `json:"usage_cpu"`
	UsageMem   float64 `json:"usage_mem"`
	Pool       string  `json:"pool"`
	Dedicated  bool    `json:"dedicated"`
	Model      string  `json:"model"`
	MaxVMs     int     `json:"max_vms"`
	CurrentVMs int     `json:"current_vms"`
}

// SchedulingRequest represents a request for scheduling a VM
type SchedulingRequest struct {
	RequestedCPU    int    `json:"requested_cpu"`
	RequestedMemory int    `json:"requested_memory"`
	RequestedPool   string `json:"requested_pool"`
	Dedicated       bool   `json:"dedicated"`
}

// TaskStatus tracks scheduling task states
type TaskStatus struct {
	Status string         `json:"status"`
	Result *BareMetalNode `json:"result,omitempty"`
}

var (
	nodesCache []BareMetalNode
	cacheLock  sync.RWMutex
	taskQueue  = make(chan schedulingTask, 100)
	taskStatus = make(map[string]*TaskStatus)
	taskLock   sync.Mutex
)

type schedulingTask struct {
	TaskID     string
	Request    SchedulingRequest
	RetryCount int
}

func fetchBareMetalNodes() []BareMetalNode {
	resp, err := http.Get("http://localhost:5000/nodes")
	if err != nil {
		log.Println("[Fetch Error]", err)
		return nil
	}
	defer resp.Body.Close()

	var nodes []BareMetalNode
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		log.Println("[Decode Error]", err)
		return nil
	}
	return nodes
}

func updateNodesCache() {
	for {
		nodes := fetchBareMetalNodes()
		cacheLock.Lock()
		nodesCache = nodes
		cacheLock.Unlock()
		log.Println("[Cache Updated]", nodesCache)
		time.Sleep(10 * time.Second)
	}
}

func preFilter(request SchedulingRequest, nodes []BareMetalNode) []BareMetalNode {
	filtered := []BareMetalNode{}
	for _, node := range nodes {
		if node.Pool == request.RequestedPool && (!request.Dedicated || node.Dedicated) && node.CurrentVMs < node.MaxVMs {
			filtered = append(filtered, node)
		}
	}
	return filtered
}

func filterNodes(request SchedulingRequest, nodes []BareMetalNode) []BareMetalNode {
	filtered := []BareMetalNode{}
	for _, node := range nodes {
		availableCPU := int((1.0 - node.UsageCPU) * float64(node.CPU))
		availableMem := int((1.0 - node.UsageMem) * float64(node.Memory))
		if availableCPU >= request.RequestedCPU && availableMem >= request.RequestedMemory {
			filtered = append(filtered, node)
		}
	}
	return filtered
}

func scoreNodes(nodes []BareMetalNode) []BareMetalNode {
	if len(nodes) == 0 {
		return nil
	}
	// Simple sorting based on available CPU & Memory
	for i := 0; i < len(nodes)-1; i++ {
		for j := i + 1; j < len(nodes); j++ {
			if (1-nodes[i].UsageCPU)+(1-nodes[i].UsageMem) < (1-nodes[j].UsageCPU)+(1-nodes[j].UsageMem) {
				nodes[i], nodes[j] = nodes[j], nodes[i]
			}
		}
	}
	return nodes
}

func selectBestNode(request SchedulingRequest, nodes []BareMetalNode) *BareMetalNode {
	preFiltered := preFilter(request, nodes)
	filtered := filterNodes(request, preFiltered)
	candidates := scoreNodes(filtered)
	if len(candidates) > 0 {
		return &candidates[0]
	}
	return nil
}

func schedulerWorker() {
	for task := range taskQueue {
		nodes := fetchBareMetalNodes()
		bestNode := selectBestNode(task.Request, nodes)

		taskLock.Lock()
		if bestNode != nil {
			taskStatus[task.TaskID] = &TaskStatus{Status: "completed", Result: bestNode}
			log.Printf("[Task %s] Scheduled to %s\n", task.TaskID, bestNode.Name)
		} else if task.RetryCount < 3 {
			delay := time.Duration(3*task.RetryCount+rand.Intn(1000)) * time.Millisecond
			log.Printf("[Retry %d] Task %s retrying in %v\n", task.RetryCount, task.TaskID, delay)
			time.Sleep(delay)
			taskQueue <- schedulingTask{TaskID: task.TaskID, Request: task.Request, RetryCount: task.RetryCount + 1}
		} else {
			taskStatus[task.TaskID] = &TaskStatus{Status: "failed"}
			log.Printf("[Task %s] Failed after retries\n", task.TaskID)
		}
		taskLock.Unlock()
	}
}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	var request SchedulingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	taskID := uuid.New().String()
	taskLock.Lock()
	taskStatus[taskID] = &TaskStatus{Status: "pending"}
	taskLock.Unlock()
	taskQueue <- schedulingTask{TaskID: taskID, Request: request, RetryCount: 0}

	json.NewEncoder(w).Encode(map[string]string{"task_id": taskID, "message": "Task submitted"})
}

func checkStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["task_id"]
	taskLock.Lock()
	status, exists := taskStatus[taskID]
	taskLock.Unlock()
	if !exists {
		http.Error(w, "Task ID not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(status)
}

func main() {
	go updateNodesCache()
	go schedulerWorker()

	r := mux.NewRouter()
	r.HandleFunc("/schedule", scheduleHandler).Methods("POST")
	r.HandleFunc("/schedule/{task_id}", checkStatusHandler).Methods("GET")

	log.Println("Scheduler API running on :8080")
	http.ListenAndServe(":8080", r)
}
