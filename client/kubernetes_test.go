package client

import (
	"os"
	"reflect"
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

func createMockKubeConfigFile(path string) (*os.File, error) {
	var data = []byte(`
apiVersion: v1
clusters:
- cluster:
    server: https://127.0.0.1:55171
  name: kind-kind
- cluster:
    server: https://127.0.0.1:55902
  name: kind-kind2
contexts:
- context:
    cluster: kind-kind
    user: kind-kind
  name: kind-kind
- context:
    cluster: kind-kind2
    user: kind-kind2
  name: kind-kind2
current-context: kind-kind2
kind: Config
preferences: {}
users:
- name: kind-kind
  user:
- name: kind-kind2
  user:
`)

	file, err := os.CreateTemp("", "temp")
	if err != nil {
		return nil, err
	}

	if _, err := file.Write(data); err != nil {
		return nil, err
	}

	return file, nil
}

func TestNewKubeClient(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *KubeClient
	}{
		{
			name: "Create a new kube client",
			args: args{
				path: "test",
			},
			want: &KubeClient{
				Path:      "test",
				Clientset: &kubernetes.Clientset{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKubeClient(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKubeClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubeClient_ClientsetWithCurrentContext(t *testing.T) {
	file, err := createMockKubeConfigFile("temp")
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	defer os.Remove(file.Name())

	type fields struct {
		Path      string
		Context   string
		Clientset kubernetes.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Create a new clientset with current context",
			fields: fields{
				Path:      file.Name(),
				Context:   "kind-kind",
				Clientset: fake.NewSimpleClientset(),
			},
			wantErr: false,
		},
		{
			name: "Try with invalid path",
			fields: fields{
				Path:      "test",
				Context:   "test",
				Clientset: fake.NewSimpleClientset(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := &KubeClient{
				Path:      tt.fields.Path,
				Context:   tt.fields.Context,
				Clientset: tt.fields.Clientset,
			}
			_, err := kube.ClientsetWithCurrentContext()
			if (err != nil) != tt.wantErr {
				t.Errorf("KubeClient.ClientsetWithCurrentContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestKubeClient_ClientSetWithContext(t *testing.T) {
	file, err := createMockKubeConfigFile("temp")
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	defer os.Remove(file.Name())
	type fields struct {
		Path      string
		Context   string
		Clientset kubernetes.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Create a new clientset with new context",
			fields: fields{
				Path:      file.Name(),
				Context:   "kind-kind2",
				Clientset: fake.NewSimpleClientset(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := &KubeClient{
				Path:      tt.fields.Path,
				Context:   tt.fields.Context,
				Clientset: tt.fields.Clientset,
			}
			_, err := kube.ClientSetWithContext()
			if (err != nil) != tt.wantErr {
				t.Errorf("KubeClient.ClientSetWithContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestKubeClient_MetricsClientSetWitContext(t *testing.T) {
	type fields struct {
		Path      string
		Context   string
		Clientset kubernetes.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		want    *versioned.Clientset
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := &KubeClient{
				Path:      tt.fields.Path,
				Context:   tt.fields.Context,
				Clientset: tt.fields.Clientset,
			}
			got, err := kube.MetricsClientSetWithContext()
			if (err != nil) != tt.wantErr {
				t.Errorf("KubeClient.MetricsClientSetWitContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KubeClient.MetricsClientSetWitContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubeClient_CreateClientSetFromS3(t *testing.T) {
	type fields struct {
		Path      string
		Context   string
		Clientset kubernetes.Interface
	}
	type args struct {
		filepath string
		awsOpt   AWSOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    kubernetes.Interface
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := &KubeClient{
				Path:      tt.fields.Path,
				Context:   tt.fields.Context,
				Clientset: tt.fields.Clientset,
			}
			got, err := kube.CreateClientSetFromS3(tt.args.filepath, tt.args.awsOpt)
			if (err != nil) != tt.wantErr {
				t.Errorf("KubeClient.CreateClientSetFromS3() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KubeClient.CreateClientSetFromS3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubeClient_NewConfigFromContext(t *testing.T) {
	type fields struct {
		Path      string
		Context   string
		Clientset kubernetes.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		want    *rest.Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := &KubeClient{
				Path:      tt.fields.Path,
				Context:   tt.fields.Context,
				Clientset: tt.fields.Clientset,
			}
			got, err := kube.NewConfigFromContext()
			if (err != nil) != tt.wantErr {
				t.Errorf("KubeClient.NewConfigFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KubeClient.NewConfigFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createClientSet(t *testing.T) {
	type args struct {
		r *rest.Config
	}
	tests := []struct {
		name    string
		args    args
		want    kubernetes.Interface
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createClientSet(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("createClientSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createClientSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubeClient_SwitchKubeContext(t *testing.T) {
	type fields struct {
		Path      string
		Context   string
		Clientset kubernetes.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := &KubeClient{
				Path:      tt.fields.Path,
				Context:   tt.fields.Context,
				Clientset: tt.fields.Clientset,
			}
			if err := kube.SwitchKubeContext(); (err != nil) != tt.wantErr {
				t.Errorf("KubeClient.SwitchKubeContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
