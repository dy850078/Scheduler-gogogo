package algorithm

import "github.com/dy850078/virtflow-scheduler-go/internal/model"

// SchedulingStrategy defines interface for pluggable node selection strategies
type SchedulingStrategy interface {
	SelectBestNode(req model.SchedulingRequest, nodes []model.BareMetalNode) *model.BareMetalNode
}
