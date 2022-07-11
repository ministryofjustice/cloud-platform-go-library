package mock

import (
	"time"

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
		m.Cluster.Nodes = v1.NodeList{
			Items: []v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "Node1",
						Labels: map[string]string{
							"Cluster": "Cluster1",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "Node2",
						CreationTimestamp: metav1.Time{
							Time: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		}
	}
}

// WithMonitoringNodes returns a MockOptions function that creates nodes with and without monitoring tags.
func WithMonitoringNodes() MockOptions {
	return func(m *Mock) {
		m.Cluster.Nodes = v1.NodeList{
			Items: []v1.Node{
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
			},
		}
	}
}

// WithBrokenNodes returns MockOptions that sets some nodes to be health and some unhealthy.
func WithBrokenNodes() MockOptions {
	return func(m *Mock) {
		m.Cluster.Nodes = v1.NodeList{
			Items: []v1.Node{
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
			},
		}
	}
}

// WithWorkingPods returns a MockOptions function that sets the pods to be running.
func WithWorkingPods() MockOptions {
	return func(m *Mock) {
		m.Cluster.Pods = &v1.PodList{
			Items: []v1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "Pod1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "Pod2",
					},
				},
			},
		}
	}
}

// WithBrokenPods returns a MockOptions function that sets some pods to be running and some not.
func WithBrokenPods() MockOptions {
	return func(m *Mock) {
		m.Cluster.Pods = &v1.PodList{
			Items: []v1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "Pod1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "Pod2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "Pod3",
					},
					Status: v1.PodStatus{
						Phase: v1.PodFailed,
					},
				},
			},
		}
	}
}

// WithNamespaces returns a MockOptions function that sets the namespaces to be ready.
func WithNamespaces() MockOptions {
	return func(m *Mock) {
		m.Cluster.Namespaces = v1.NamespaceList{
			Items: []v1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "Namespace1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "Namespace2",
					},
				},
			},
		}
	}
}
