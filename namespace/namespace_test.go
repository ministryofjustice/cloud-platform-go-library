package namespace

import (
	"reflect"
	"testing"

	"github.com/ministryofjustice/cloud-platform-go-library/client"
	v1 "k8s.io/api/core/v1"
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
		// TODO: Add test cases.
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
