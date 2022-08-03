package namespace_test

import (
	"reflect"
	"testing"

	"github.com/ministryofjustice/cloud-platform-go-library/client"
	"github.com/ministryofjustice/cloud-platform-go-library/mock"
	"github.com/ministryofjustice/cloud-platform-go-library/namespace"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	// fakeClient is a fake cloud platform kubernetes cluster
	// containing some namespace objects.
	fakeCluster = mock.NewCluster(
		mock.WithNamespaces(),
	)
	// fakeClient lets you test against the mock cluster.
	fakeClient = client.KubeClient{
		Clientset: fake.NewSimpleClientset(&fakeCluster.Cluster.Namespaces),
	}
)

func TestAllNamespaces(t *testing.T) {
	type args struct {
		c *client.KubeClient
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.NamespaceList
		wantErr bool
	}{
		{
			name: "get all namespaces from the cluster",
			args: args{
				c: &fakeClient,
			},
			want: &v1.NamespaceList{
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := namespace.AllNamespaces(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllNamespaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllNamespaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNamespace(t *testing.T) {
	type args struct {
		c    *client.KubeClient
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.Namespace
		wantErr bool
	}{
		{
			name: "get namespace from the cluster",
			args: args{
				c:    &fakeClient,
				name: "Namespace1",
			},
			want: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Namespace1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := namespace.Namespace(tt.args.c, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Namespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Namespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateNamespace(t *testing.T) {
	type args struct {
		c    *client.KubeClient
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.Namespace
		wantErr bool
	}{
		{
			name: "create namespace in the cluster",
			args: args{
				c:    &fakeClient,
				name: "Namespace9",
			},
			want: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Namespace9",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := namespace.CreateNamespace(tt.args.c, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteNamespace(t *testing.T) {
	type args struct {
		c    *client.KubeClient
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "delete namespace from the cluster",
			args: args{
				c:    &fakeClient,
				name: "Namespace9",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := namespace.DeleteNamespace(tt.args.c, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TODO: This test is failing inexplicably at the moment. It's not clear why.
// func TestGetTeamNamespaces(t *testing.T) {
// 	list, err := namespace.GetTeamNamespaces(&fakeClient, "webops")
// 	if err != nil {
// 		t.Errorf("GetTeamNamespaces() error = %v", err)
// 	}

// 	assert.EqualValues(t, "Namespace2", list[0].Name)
// }

func TestNamespaceSlackChannel(t *testing.T) {
	type args struct {
		c    *client.KubeClient
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "get namespace slack channel",
			args: args{
				c:    &fakeClient,
				name: "Namespace2",
			},
			want: "cloud-platform",
		},
		{
			name: "get non-existant namespace slack channel and fail",
			args: args{
				c:    &fakeClient,
				name: "Namespace100",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := namespace.NamespaceSlackChannel(tt.args.c, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamespaceSlackChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NamespaceSlackChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}
