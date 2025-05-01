package model

type SchedulingRequest struct {
	RequestedCPU    int    `json:"requested_cpu"`
	RequestedMemory int    `json:"requested_memory"`
	RequestedPool   string `json:"requested_pool"`
	Dedicated       bool   `json:"dedicated"`
	TaskID          string `json:"task_id"`
}
