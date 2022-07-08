package cluster_test

import (
	"reflect"
	"testing"

	"github.com/ministryofjustice/cloud-platform-go-library/client"
	"github.com/ministryofjustice/cloud-platform-go-library/cluster"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestAllPods(t *testing.T) {
	type args struct {
		c client.KubeClient
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.PodList
		wantErr bool
	}{
		{
			name: "get all pods with working client",
			args: args{
				c: client.KubeClient{
					Clientset: fake.NewSimpleClientset(standard.Cluster.Pods),
				},
			},
			want: &v1.PodList{
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
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cluster.AllPods(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllPods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllPods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStuckPods(t *testing.T) {
	client := client.KubeClient{
		Clientset: fake.NewSimpleClientset(broken.Cluster.Pods),
	}

	pods, err := cluster.AllPods(client)
	if err != nil {
		t.Errorf("AllPods() error = %v", err)
	}

	stuckPods, err := cluster.StuckPods(client, *pods)
	if err != nil {
		t.Errorf("StuckPods() error = %v", err)
	}

	assert.Equal(t, 1, len(stuckPods))
}
