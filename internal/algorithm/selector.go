
package algorithm

import "virtflow-go/internal/model"

func SelectBestNode(req model.SchedulingRequest, nodes []model.BareMetalNode) *model.BareMetalNode {
    var selected *model.BareMetalNode
    for _, node := range nodes {
        if node.Pool != req.RequestedPool || node.Dedicated != req.Dedicated {
            continue
        }
        if node.CPU < req.RequestedCPU || node.Memory < req.RequestedMemory {
            continue
        }
        if selected == nil || node.UsageCPU < selected.UsageCPU {
            selected = &node
        }
    }
    return selected
}
