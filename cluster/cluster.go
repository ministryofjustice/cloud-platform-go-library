package cluster

import (
	"github.com/ministryofjustice/cloud-platform-go-library/client"
	v1 "k8s.io/api/core/v1"
)

// Cluster struct represents an MoJ Cloud Platform Kubernetes cluster object
type Cluster struct {
	Name       string
	NewestNode v1.Node
	Nodes      v1.NodeList
	OldestNode v1.Node
	Pods       *v1.PodList
	StuckPods  []*v1.Pod
}

// NewWithValues returns a full Cluster object with populated values.
func NewWithValues(c client.KubeClient) (*Cluster, error) {
	nodes, err := AllNodes(c)
	if err != nil {
		return nil, err
	}
	pods, err := AllPods(c)
	if err != nil {
		return nil, err
	}
	stuckPods, err := StuckPods(c, *pods)
	if err != nil {
		return nil, err
	}
	cluster := &Cluster{
		Nodes:     nodes,
		Pods:      pods,
		StuckPods: stuckPods,
	}

	// You can only get the name of a Cloud Platform cluster using the labels on a node.
	cluster.GetName()

	return cluster, nil
}

// GetName is a method function to get the name of the cluster.
func (c *Cluster) GetName() {
	c.Name = c.Nodes.Items[0].Labels["Cluster"]
	return
}
