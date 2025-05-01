package algorithm

import "github.com/dy850078/virtflow-scheduler-go/internal/model"

// TODO: implement hybrid CPU+MEM scoring strategy
type HybridStrategy struct{}

func (s *HybridStrategy) SelectBestNode(req model.SchedulingRequest, nodes []model.BareMetalNode) *model.BareMetalNode {
	return nil
}
