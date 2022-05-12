package client

import (
	"log"
	"os"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

var file *os.File

func TestMain(m *testing.M) {
	var err error
	file, err = createMockKubeConfigFile("temp")
	if err != nil {
		log.Fatalln(err)
	}
	code := m.Run()

	defer file.Close()
	defer os.Remove(file.Name())
	os.Exit(code)
}

func TestNewKubeClientWithValues(t *testing.T) {
	type args struct {
		configFilePath string
		context        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create new client with empty context",
			args: args{
				configFilePath: file.Name(),
				context:        "",
			},
			wantErr: false,
		},
		{
			name: "Create new client with existing context",
			args: args{
				configFilePath: file.Name(),
				context:        "kind-kind",
			},
			wantErr: false,
		},
		{
			name: "Create new client with fake context",
			args: args{
				configFilePath: file.Name(),
				context:        "fake",
			},
			wantErr: true,
		},
		{
			name: "Create new client with incorrect file path",
			args: args{
				configFilePath: "fake",
				context:        "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewKubeClientWithValues(tt.args.configFilePath, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKubeClientWithValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMetricsKubeClientWithValues(tt.args.configFilePath, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetricsKubeClientWithValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestKubeClient_BuildClientSet(t *testing.T) {
	type fields struct {
		Path               string
		Context            string
		Clientset          kubernetes.Interface
		VersionedClientset versioned.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Build clientset with correct data",
			fields: fields{
				Path:    file.Name(),
				Context: "kind-kind",
			},
			wantErr: false,
		},
		{
			name: "Build clientset with incorrect data",
			fields: fields{
				Path:    file.Name(),
				Context: "fake",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := &KubeClient{
				Path:               tt.fields.Path,
				Context:            tt.fields.Context,
				Clientset:          tt.fields.Clientset,
				VersionedClientset: tt.fields.VersionedClientset,
			}
			if err := kube.BuildClientSet(); (err != nil) != tt.wantErr {
				t.Errorf("KubeClient.BuildClientSet() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := kube.BuildVersionedClientset(); (err != nil) != tt.wantErr {
				t.Errorf("KubeClient.VersionedClientset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKubeClient_NewConfigFromContext(t *testing.T) {
	type fields struct {
		Path               string
		Context            string
		Clientset          kubernetes.Interface
		VersionedClientset versioned.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		want    *rest.Config
		wantErr bool
	}{
		{
			name: "Create new config",
			fields: fields{
				Path:    file.Name(),
				Context: "kind-kind",
			},
			want: &rest.Config{
				Host:    "https://127.0.0.1:55171",
				APIPath: "", ContentConfig: rest.ContentConfig{
					AcceptContentTypes:   "",
					ContentType:          "",
					GroupVersion:         (*schema.GroupVersion)(nil),
					NegotiatedSerializer: runtime.NegotiatedSerializer(nil),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := &KubeClient{
				Path:               tt.fields.Path,
				Context:            tt.fields.Context,
				Clientset:          tt.fields.Clientset,
				VersionedClientset: tt.fields.VersionedClientset,
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

func TestKubeClient_SwitchKubeContext(t *testing.T) {
	type fields struct {
		Path               string
		Context            string
		Clientset          kubernetes.Interface
		VersionedClientset versioned.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Switch the kube context using correct data",
			fields: fields{
				Path:    file.Name(),
				Context: "kind-kind",
			},
			wantErr: false,
		},
		{
			name: "Switch the kube context using incorrect data",
			fields: fields{
				Path:    file.Name(),
				Context: "fake",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := &KubeClient{
				Path:               tt.fields.Path,
				Context:            tt.fields.Context,
				Clientset:          tt.fields.Clientset,
				VersionedClientset: tt.fields.VersionedClientset,
			}
			if err := kube.SwitchKubeContext(); (err != nil) != tt.wantErr {
				t.Errorf("KubeClient.SwitchKubeContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
