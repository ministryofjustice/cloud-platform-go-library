package namespace

import (
	"os"
	"reflect"
	"testing"

	"github.com/ministryofjustice/cloud-platform-go-library/client"
	"github.com/ministryofjustice/cloud-platform-go-library/mock"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	fakeCluster *mock.Mock
	fakeClient  *client.KubeClient
)

func TestMain(m *testing.M) {
	fakeCluster = mock.NewCluster(
		mock.WithNamespaces(),
	)
	fakeClient.Clientset = fake.NewSimpleClientset(&fakeCluster.Cluster.Namespaces)
	code := m.Run()
	os.Exit(code)
}

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
				c: fakeClient,
			},
			want:    &fakeCluster.Cluster.Namespaces,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AllNamespaces(tt.args.c)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Namespace(tt.args.c, tt.args.name)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateNamespace(tt.args.c, tt.args.name)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteNamespace(tt.args.c, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTeamNamespaces(t *testing.T) {
	type args struct {
		team string
	}
	tests := []struct {
		name    string
		args    args
		want    []*v1.Namespace
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTeamNamespaces(tt.args.team)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTeamNamespaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTeamNamespaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NamespaceSlackChannel(tt.args.c, tt.args.name)
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

func TestProductionNamespace(t *testing.T) {
	type args struct {
		c *client.KubeClient
	}
	tests := []struct {
		name    string
		args    args
		want    []*v1.Namespace
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProductionNamespace(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductionNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductionNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNonProductionNamespace(t *testing.T) {
	type args struct {
		c *client.KubeClient
	}
	tests := []struct {
		name    string
		args    args
		want    []*v1.Namespace
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NonProductionNamespace(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NonProductionNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NonProductionNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNamespaceSourceCode(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NamespaceSourceCode(tt.args.c, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamespaceSourceCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NamespaceSourceCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNamespaceOwner(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NamespaceOwner(tt.args.c, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NamespaceOwner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NamespaceOwner() = %v, want %v", got, tt.want)
			}
		})
	}
}
