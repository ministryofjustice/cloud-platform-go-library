package utils

import (
	"testing"
)

func TestGetOwnerRepoPull(t *testing.T) {
	type args struct {
		ref  string
		repo string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 int
	}{
		{
			name: "Test GetOwnerRepoPull",
			args: args{
				ref:  "refs/pull/1/merge",
				repo: "ministryofjustice/cloud-platform-environments",
			},
			want:  "ministryofjustice",
			want1: "cloud-platform-environments",
			want2: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := GetOwnerRepoPull(tt.args.ref, tt.args.repo)
			if got != tt.want {
				t.Errorf("GetOwnerRepoPull() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetOwnerRepoPull() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("GetOwnerRepoPull() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != nil {
				t.Errorf("GetOwnerRepoPull() got3 = %v, want %v", got3, nil)
			}
		})
	}
}

func TestValidateModuleSource(t *testing.T) {
	type args struct {
		source          string
		approvedModules map[string]bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Test ValidateModuleSource True Multiple Modules",
			args: args{
				source: "github.com/ministryofjustice/cloud-platform-terraform-ecr-credentials",
				approvedModules: map[string]bool{
					"github.com/ministryofjustice/cloud-platform-terraform-rds-instance":     false,
					"github.com/ministryofjustice/cloud-platform-terraform-rds-aurora":       false,
					"github.com/ministryofjustice/cloud-platform-terraform-serviceaccount":   false,
					"github.com/ministryofjustice/cloud-platform-terraform-dynamodb-cluster": false,
					"github.com/ministryofjustice/cloud-platform-terraform-sqs":              false,
					"github.com/ministryofjustice/cloud-platform-terraform-s3-bucket":        false,
					"github.com/ministryofjustice/cloud-platform-terraform-sns-topic":        false,
					"github.com/ministryofjustice/cloud-platform-terraform-secrets-manager":  false,
					"github.com/ministryofjustice/cloud-platform-terraform-opensearch":       false,
					"github.com/ministryofjustice/cloud-platform-terraform-ecr-credentials":  true,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test ValidateModuleSource False Multiple Modules",
			args: args{
				source: "github.com/ministryofjustice/cloud-platform-terraform-ecr-credentials",
				approvedModules: map[string]bool{
					"github.com/ministryofjustice/cloud-platform-terraform-rds-instance":     true,
					"github.com/ministryofjustice/cloud-platform-terraform-rds-aurora":       false,
					"github.com/ministryofjustice/cloud-platform-terraform-serviceaccount":   false,
					"github.com/ministryofjustice/cloud-platform-terraform-dynamodb-cluster": false,
					"github.com/ministryofjustice/cloud-platform-terraform-sqs":              false,
					"github.com/ministryofjustice/cloud-platform-terraform-s3-bucket":        false,
					"github.com/ministryofjustice/cloud-platform-terraform-sns-topic":        false,
					"github.com/ministryofjustice/cloud-platform-terraform-secrets-manager":  false,
					"github.com/ministryofjustice/cloud-platform-terraform-opensearch":       false,
					"github.com/ministryofjustice/cloud-platform-terraform-ecr-credentials":  false,
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateModuleSource(tt.args.source, tt.args.approvedModules)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateModuleSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateModuleSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
