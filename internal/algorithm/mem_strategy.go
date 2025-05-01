package algorithm

import "github.com/dy850078/virtflow-scheduler-go/internal/model"

// TODO: implement memory-based node selector
type MemoryStrategy struct{}

func (s *MemoryStrategy) SelectBestNode(req model.SchedulingRequest, nodes []model.BareMetalNode) *model.BareMetalNode {
	return nil
}
