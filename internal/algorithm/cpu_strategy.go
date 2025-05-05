package algorithm

import (
	"sort"

	"github.com/dy850078/virtflow-scheduler-go/internal/model"
)

// CPUStrategy selects node with lowest CPU usage that meets the request
type CPUStrategy struct{}

func (s *CPUStrategy) SelectBestNode(req model.SchedulingRequest, nodes []model.BareMetalNode) *model.BareMetalNode {
	var candidates []model.BareMetalNode

	for _, node := range nodes {
		if node.Pool != req.RequestedPool || node.Dedicated != req.Dedicated {
			continue
		}
		if node.CPU < req.RequestedCPU || node.Memory < req.RequestedMemory {
			continue
		}
		candidates = append(candidates, node)
	}

	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].UsageCPU < candidates[j].UsageCPU
	})

	return &candidates[0]
}
