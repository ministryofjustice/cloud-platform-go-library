package cluster

import (
	"reflect"
	"testing"

	"github.com/ministryofjustice/cloud-platform-go-library/client"
)

func TestNewWithValues(t *testing.T) {
	type args struct {
		c *client.KubeClient
	}
	tests := []struct {
		name    string
		args    args
		want    *Cluster
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWithValues(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWithValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWithValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
