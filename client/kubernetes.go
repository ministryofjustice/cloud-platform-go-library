package client

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// KubeClient is used to pass kubeconfig options
// to methods that need to interact with the kubernetes api
type KubeClient struct {
	// ConfigPath is an optional filepath to a kubeconfig file.
	// Usually this is set to ~/.kube/config.
	Path string
	// Context is the name of the context to use in the kubeconfig file.
	Context string
	// Clientset is the kubernetes client set required to interact with the kubernetes api
	Clientset kubernetes.Interface
}

// AWSOptions is used to pass aws options to methods that need them
type AWSOptions struct {
	Profile string
	Key     string
	Secret  string
	Region  string
	Bucket  string
}

func NewKubeClient(path string) *KubeClient {
	if path == "" {
		path = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	return &KubeClient{
		Path:      path,
		Clientset: &kubernetes.Clientset{},
	}
}

func (kube *KubeClient) ClientsetWithCurrentContext() (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kube.Path)
	if err != nil {
		return nil, err
	}

	return createClientSet(config)
}

func (kube *KubeClient) ClientSetWithContext() (kubernetes.Interface, error) {
	client, err := kube.NewConfigFromContext()
	if err != nil {
		return nil, err
	}

	return createClientSet(client)
}

func (kube *KubeClient) MetricsClientSetWithContext() (*versioned.Clientset, error) {
	client, err := kube.NewConfigFromContext()
	if err != nil {
		return nil, err
	}

	return versioned.NewForConfig(client)
}

func (kube *KubeClient) CreateClientSetFromS3(filepath string, awsOpt AWSOptions) (kubernetes.Interface, error) {
	// Ignore the kubeconfig filepath variable and download the config from s3
	_, err := kube.DownloadS3Kubeconfig(filepath, awsOpt)
	if err != nil {
		return nil, err
	}

	client, err := kube.ClientsetWithCurrentContext()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (kube *KubeClient) NewConfigFromContext() (*rest.Config, error) {
	if kube.Path == "" || kube.Context == "" {
		return nil, fmt.Errorf("Config file and Context must be set")
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kube.Path},
		&clientcmd.ConfigOverrides{
			CurrentContext: kube.Context,
		}).ClientConfig()
}

func createClientSet(r *rest.Config) (kubernetes.Interface, error) {
	return kubernetes.NewForConfig(r)
}

func (kube *KubeClient) SwitchKubeContext() error {
	kubeconf := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kube.Path},
		&clientcmd.ConfigOverrides{
			CurrentContext: kube.Context,
		})

	config, err := kubeconf.RawConfig()
	if err != nil {
		return err
	}

	if config.Contexts[kube.Context] == nil {
		return fmt.Errorf("context %s not found", kube.Context)
	}

	config.CurrentContext = kube.Context
	err = clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), config, true)
	if err != nil {
		return err
	}

	return nil
}

func (kube *KubeClient) DownloadS3Kubeconfig(fileName string, awsOpt AWSOptions) (string, error) {
	if awsOpt.Key == "" && awsOpt.Secret == "" || awsOpt.Profile == "" {
		return "", fmt.Errorf("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY or AWS_PROFILE must be set")
	}

	if awsOpt.Region == "" {
		awsOpt.Region = "eu-west-2"
	}

	buff := &aws.WriteAtBuffer{}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsOpt.Region),
	})
	if err != nil {
		return "", err
	}

	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(awsOpt.Bucket),
			Key:    aws.String(fileName),
		})
	if err != nil {
		return "", err
	}
	if numBytes < 1 {
		return "", fmt.Errorf("error the kubecfg file downloaded is empty and must have failed")
	}

	data := buff.Bytes()
	err = ioutil.WriteFile(kube.Path, data, 0o644)
	if err != nil {
		return "", err
	}

	return kube.Path, nil
}
