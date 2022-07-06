package cluster

import (
	"context"
	"fmt"

	"github.com/ministryofjustice/cloud-platform-go-library/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AllNodes returns all nodes in a cluster (including master nodes) as a slice of v1.Node objects.
func AllNodes(c client.KubeClient) ([]v1.Node, error) {
	nodeList, err := c.Clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %s", err)
	}

	return nodeList.Items, nil
}

// MonitoringNodes returns all nodes in the cluster tagged as monitoring nodes.
func MonitoringNodes(c client.KubeClient) ([]*v1.Node, error) {
	nodes, err := AllNodes(c)
	if err != nil {
		return nil, err
	}

	var monitoringNodes []*v1.Node
	for _, node := range nodes {
		if node.Labels["monitoring_ng"] == "true" {
			monitoringNodes = append(monitoringNodes, &node)
		}
	}
	return monitoringNodes, nil
}

// OldestNode returns the oldest node in a slice of v1.Node objects.
// It uses the creation timestamp to determine the oldest node.
func OldestNode(c client.KubeClient, nodes []*v1.Node) (*v1.Node, error) {
	oldest := nodes[0]

	for _, node := range nodes {
		if node.CreationTimestamp.Before(&oldest.CreationTimestamp) && node.Spec.Taints == nil {
			oldest = node
		}
	}

	if oldest == nil {
		return nil, fmt.Errorf("oldest node not found, are you sure you passed a full slice of nodes?")
	}

	return oldest, nil
}

// NewestNode returns the newest node in a slice of v1.Node objects.
// It uses the creation timestamp to determine the newest node.
func NewestNode(c client.KubeClient, nodes []v1.Node) v1.Node {
	newest := nodes[0]

	// Tried to use the After method on the creation timestamp, but it requires time.Time objects,
	// not v1.Node objects.
	for _, node := range nodes {
		date := node.CreationTimestamp.Time
		if date.After(newest.CreationTimestamp.Time) {
			newest = node
		}
	}

	return newest
}

