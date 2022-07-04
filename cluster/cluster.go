package cluster

import (
	"github.com/ministryofjustice/cloud-platform-go-library/client"
	v1 "k8s.io/api/core/v1"
)

// Clustr struct represents an MoJ Cloud Platform Kubernetes cluster object
type Cluster struct {
	Name       string
	NewestNode v1.Node
	Nodes      []v1.Node
	OldestNode v1.Node
	Pods       []v1.Pod
	StuckPods  []v1.Pod
}

// NewWithValues returns a full Cluster object with populated values.
func NewWithValues(c *client.KubeClient) (*Cluster, error) {
	return nil, nil
}

func (c *Cluster) GetName() string {
	return c.Nodes[0].Labels["Cluster"]
}
