package mock

import (
	"github.com/ministryofjustice/cloud-platform-go-library/cluster"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Mock is a mock implementation of the Cluster type defined in the cluster package.
type Mock struct {
	Cluster cluster.Cluster
}

// MockOptions is a function type used to modify the Mock object.
type MockOptions func(*Mock)

// NewCluster returns a new Mock object.
func NewCluster(opts ...MockOptions) *Mock {
	m := &Mock{
		Cluster: cluster.Cluster{},
	}

	for _, opt := range opts {
		opt(m)
	}
	return m
}

// WithWorkingNodes returns a MockOptions function that sets the nodes to be ready and healthy.
func WithWorkingNodes() MockOptions {
	return func(m *Mock) {
		m.Cluster.Nodes = []v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Node1",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Node2",
				},
			},
		}
	}
}

// WithMonitoringNodes returns a MockOptions function that creates nodes with and without monitoring tags.
func WithMonitoringNodes() MockOptions {
	return func(m *Mock) {
		m.Cluster.Nodes = []v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Node1",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Node2",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Node3",
					Labels: map[string]string{
						"monitoring_ng": "true",
					},
				},
			},
		}
	}
}

// WithBrokenNodes returns MockOptions that sets some nodes to be health and some unhealthy.
func WithBrokenNodes() MockOptions {
	return func(m *Mock) {
		m.Cluster.Nodes = []v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Node1",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Node2",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Node3",
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:   v1.NodeDiskPressure,
							Status: v1.ConditionUnknown,
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Node4",
				},
				Status: v1.NodeStatus{
					Conditions: []v1.NodeCondition{
						{
							Type:   v1.NodeReady,
							Status: v1.ConditionFalse,
						},
					},
				},
			},
		}
	}
}
