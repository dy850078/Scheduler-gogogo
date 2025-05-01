package model

type SchedulingRequest struct {
	RequestedCPU    int    `json:"requested_cpu"`
	RequestedMemory int    `json:"requested_memory"`
	RequestedPool   string `json:"requested_pool"`
	Dedicated       bool   `json:"dedicated"`
	TaskID          string `json:"task_id"`
}
type BareMetalNode struct {
	Name        string
	CPU         int
	Memory      int
	Pool        string
	Dedicated   bool
	UsageCPU    float64
	UsageMemory float64
}

func MockNodes() []BareMetalNode {
	return []BareMetalNode{
		{"bm01", 16, 32768, "default", false, 0.3, 0.4},
		{"bm02", 32, 65536, "high-performance", true, 0.7, 0.6},
	}
}
