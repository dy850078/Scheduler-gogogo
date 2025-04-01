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

type BareMetalNode struct {
	Name        string  `json: "name"`
	CPU         int     `json: "cpu"`
	Memory      int     `json: "memory"`
	Storage     int     `json: "storage"`
	UsageCPU    float64 `json: "usage_cpu"`
	UsageMem    float64 `json: "usage_mem"`
	Pool        string  `json: "pool"`
	Dedicated   bool    `json: "dedicated"`
	Model       string  `json: "model"`
	Max_VM      int     `json: "max_vm"`
	Current_VMs int     `json: "current_vm"`
}

type SchedulingRequest struct {
	RequestedCPU  int    `json: "requested_cpu"`
	RequestedMem  int    `json: "requested_mem"`
	RequestedPool string `json "requested_pool"`
	Dedicated     bool   `json: "dedicated"`
}

type TaskStatus struct {
	Status string         `json: "status"`
	Result *BareMetalNode `json: "result"`
}

type schedulingTask struct {
	TaskID     string
	Request    SchedulingRequest
	RetryCount int
}

var (
	nodesCache []BareMetalNode
	cacheLock  sync.RWMutex
	taskQueue  = make(chan schedulingTask, 100)
	taskStatus = make(map[string]*TaskStatus)
	taskLock   sync.Mutex
)

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

func updateNodeCache() {
	for {
		cacheLock.Lock()
		nodesCache = fetchBareMetalNodes()
		cacheLock.Unlock()
		log.Println("[Cache Updated]", nodesCache)
		time.Sleep(10 * time.Second)
	}
}

func preFilterNodes(request SchedulingRequest, nodes []BareMetalNode) []BareMetalNode {
	preFiltered := []BareMetalNode{}
	for _, node := range nodes {
		if node.Pool == request.RequestedPool && (!request.Dedicated || node.Dedicated) && node.Current_VMs < node.Max_VM {
			preFiltered = append(preFiltered, node)
		}
	}
	return preFiltered
}

func filterNodes(request SchedulingRequest, nodes []BareMetalNode) []BareMetalNode {
	filtered := []BareMetalNode{}
	for _, node := range nodes {
		availableCPU := int((1.0 - node.UsageCPU) * float64(node.CPU))
		availableMem := int((1.0 - node.UsageMem) * float64(node.Memory))

		if availableCPU >= request.RequestedCPU && availableMem >= request.RequestedMem {
			filtered = append(filtered, node)
		}
	}
	return filtered
}

func scoreNodes([]BareMetalNode) *BareMetalNode {

}

func selectBestNode(request SchedulingRequest, nodes []BareMetalNode) *BareMetalNode {
	preFiltered := preFilterNodes(request, nodes)
	filtered := filterNodes(request, preFiltered)
	candidates := scoreNodes(filtered)
	if len(candidates) {
		return &candidates[0]
	}
	return nil
}

func scheduleWorker() {
	for task := range taskQueue {
		cacheLock.RLock()
		nodes := nodesCache
		cacheLock.RUnlock()

		if len(nodes) == 0 {
			log.Println("[Scheduler] Cache is empty, fetching nodes from API")
			nodes = fetchBareMetalNodes()
		}

		bestNode := selectBestNode(task.Request, nodes)

		taskLock.Lock()
		if bestNode != nil {
			taskStatus[task.TaskID] = &TaskStatus{Status: "completed", Result: bestNode}
			log.Printf("[Task %s] Scheduled to %s \n", task.TaskID, bestNode.Name)
		} else if task.RetryCount < 3 {
			delay := time.Duration(3*task.RetryCount+rand.Intn(10)) * time.Second
			log.Printf("[Retry %d] Task %s retrying in %d\n", task.RetryCount, task.TaskID, delay)
			time.Sleep(delay)
			taskQueue <- schedulingTask{TaskID: task.TaskID, Request: task.Request, RetryCount: task.RetryCount + 1}
		} else {
			taskStatus[task.TaskID] = &TaskStatus{Status: "failed"}
			log.Printf("[Task %s] Failed after %d retries\n", task.TaskID, task.RetryCount)
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/schedule", scheduleHandler).Methods("POST")
	// r.HandleFunc("/schedule/{task_id}", checkStatusHandler).Methods("GET")

	log.Println("Scheduler API running on :8080")
	http.ListenAndServe(":8080", r)
}
