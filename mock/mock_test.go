package mock

import (
	"reflect"
	"testing"
	"time"

	"github.com/ministryofjustice/cloud-platform-go-library/cluster"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewCluster(t *testing.T) {
	type args struct {
		opts []MockOptions
	}
	tests := []struct {
		name string
		args args
		want *Mock
	}{
		{
			name: "Create mock cluster with working nodes",
			args: args{
				opts: []MockOptions{
					WithWorkingNodes(),
				},
			},
			want: &Mock{
				cluster.Cluster{
					Nodes: []v1.Node{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "Node1",
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
				},
			},
		},
		{
			name: "Create a mock cluster with monitoring nodes enabled",
			args: args{
				opts: []MockOptions{
					WithMonitoringNodes(),
				},
			},
			want: &Mock{
				cluster.Cluster{
					Nodes: []v1.Node{
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
				},
			},
		},
		{
			name: "Create a mock cluster with broken nodes",
			args: args{
				opts: []MockOptions{
					WithBrokenNodes(),
				},
			},
			want: &Mock{
				cluster.Cluster{
					Nodes: []v1.Node{
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
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCluster(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}
