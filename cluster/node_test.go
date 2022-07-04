package cluster_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ministryofjustice/cloud-platform-go-library/client"
	"github.com/ministryofjustice/cloud-platform-go-library/cluster"
	"github.com/ministryofjustice/cloud-platform-go-library/mock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// standard and monitoring represent mock clusters with working and monitoring nodes
var standard, monitoring *mock.Mock

func TestMain(m *testing.M) {
	standard = mock.NewCluster(mock.WithWorkingNodes())
	monitoring = mock.NewCluster(mock.WithMonitoringNodes())
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
		want    []v1.Node
		wantErr bool
	}{
		{
			name: "get all nodes with working client",
			args: args{
				c: client.KubeClient{
					Clientset: fake.NewSimpleClientset(&standard.Cluster.Nodes[0]),
				},
			},
			want: []v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "Node1",
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
					Clientset: fake.NewSimpleClientset(&monitoring.Cluster.Nodes[2]),
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
	type args struct {
		c     client.KubeClient
		nodes []v1.Node
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.Node
		wantErr bool
	}{
		{
			name: "Get newest node",
			args: args{
				c: client.KubeClient{
					Clientset: fake.NewSimpleClientset(&standard.Cluster.Nodes[0]),
				},
				nodes: standard.Cluster.Nodes,
			},
			want: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Node2",
					CreationTimestamp: metav1.Time{
						Time: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cluster.NewestNode(tt.args.c, tt.args.nodes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewestNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
