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
					Nodes: v1.NodeList{
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
					Nodes: v1.NodeList{
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
					Nodes: v1.NodeList{
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
					},
				},
			},
		},
		{
			name: "Create a mock cluster with working pods",
			args: args{
				opts: []MockOptions{
					WithWorkingPods(),
				},
			},

			want: &Mock{
				cluster.Cluster{
					Pods: &v1.PodList{
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
				},
			},
		},
		{
			name: "Create a mock cluster with namespaces",
			args: args{
				opts: []MockOptions{
					WithNamespaces(),
				},
			},
			want: &Mock{
				cluster.Cluster{
					Namespaces: v1.NamespaceList{
						Items: []v1.Namespace{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "Namespace1",
								},
							},
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "Namespace2",
									Labels: map[string]string{
										"cloud-platform.justice.gov.uk/is-production":    "true",
										"cloud-platform.justice.gov.uk/environment-name": "production",
									},
									Annotations: map[string]string{
										"cloud-platform.justice.gov.uk/business-unit": "HQ",
										"cloud-platform.justice.gov.uk/slack-channel": "cloud-platform",
										"cloud-platform.justice.gov.uk/application":   "Namespace to test Terraform resources",
										"cloud-platform.justice.gov.uk/owner":         "Cloud Platform: platforms@digital.justice.gov.uk",
										"cloud-platform.justice.gov.uk/source-code":   "https://github.com/ministryofjustice/cloud-platform",
										"cloud-platform.justice.gov.uk/team-name":     "webops",
										"cloud-platform.justice.gov.uk/review-after":  "12.12.2019",
									},
								},
							},
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "Namespace3",
									Labels: map[string]string{
										"cloud-platform.justice.gov.uk/is-production":    "false",
										"cloud-platform.justice.gov.uk/environment-name": "development",
									},
									Annotations: map[string]string{
										"cloud-platform.justice.gov.uk/business-unit": "HMPPS",
										"cloud-platform.justice.gov.uk/slack-channel": "fake-channel",
										"cloud-platform.justice.gov.uk/application":   "Really cool app",

										"cloud-platform.justice.gov.uk/owner":        "Really cool team",
										"cloud-platform.justice.gov.uk/source-code":  "https://github.com/ministryofjustice/not-cloud-platform",
										"cloud-platform.justice.gov.uk/team-name":    "noops",
										"cloud-platform.justice.gov.uk/review-after": "12.11.2019",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Create a mock cluster with broken pods",
			args: args{
				opts: []MockOptions{
					WithBrokenPods(),
				},
			},
			want: &Mock{
				cluster.Cluster{
					Pods: &v1.PodList{
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
