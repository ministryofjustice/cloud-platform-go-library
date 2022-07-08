package cluster_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ministryofjustice/cloud-platform-go-library/client"
	"github.com/ministryofjustice/cloud-platform-go-library/cluster"
	"github.com/ministryofjustice/cloud-platform-go-library/mock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// standard, monitoring and broken represent mock clusters with working, monitoring and broken components
var standard, monitoring, broken *mock.Mock

func TestMain(m *testing.M) {
	standard = mock.NewCluster(
		mock.WithWorkingNodes(),
		mock.WithWorkingPods(),
	)
	monitoring = mock.NewCluster(mock.WithMonitoringNodes())
	broken = mock.NewCluster(
		mock.WithBrokenNodes(),
		mock.WithBrokenPods(),
	)

	code := m.Run()
	os.Exit(code)
}

func TestAllNodes(t *testing.T) {
	type args struct {
		c client.KubeClient
	}
	tests := []struct {
		name    string
		args    args
		want    v1.NodeList
		wantErr bool
	}{
		{
			name: "get all nodes with working client",
			args: args{
				c: client.KubeClient{
					Clientset: fake.NewSimpleClientset(&standard.Cluster.Nodes),
				},
			},
			want: v1.NodeList{
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
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cluster.AllNodes(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMonitoringNodes(t *testing.T) {
	type args struct {
		c client.KubeClient
	}
	tests := []struct {
		name    string
		args    args
		want    []*v1.Node
		wantErr bool
	}{
		{
			name: "Get monitoring nodes",
			args: args{
				c: client.KubeClient{
					Clientset: fake.NewSimpleClientset(&monitoring.Cluster.Nodes),
				},
			},
			want: []*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "Node3",
						Labels: map[string]string{
							"monitoring_ng": "true",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cluster.MonitoringNodes(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("MonitoringNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MonitoringNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewestNode(t *testing.T) {
	mockClient := client.KubeClient{
		Clientset: fake.NewSimpleClientset(&monitoring.Cluster.Nodes),
	}

	nodes, err := cluster.AllNodes(mockClient)
	if err != nil {
		t.Errorf("AllNodes() error = %v", err)
	}

	assert.Equal(t, "Node1", cluster.NewestNode(mockClient, nodes.Items).Name)
}
