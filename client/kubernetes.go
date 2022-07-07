package client

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// KubeClient is used to pass kubeconfig options
// to methods that need to interact with the kubernetes api
type KubeClient struct {
	// Path is an optional filepath to a kubeconfig file.
	// Usually this is set to ~/.kube/config.
	Path string
	// Context is the name of the context to use in the kubeconfig file.
	Context string
	// Clientset is the kubernetes client set required to interact with the kubernetes api.
	Clientset kubernetes.Interface
	// VersionedClientset allows you to communicate with the kubernetes api to get metrics data.
	VersionedClientset versioned.Interface
}

// AwsOptions is used to pass aws options to functions/methods that need them
type AwsOptions struct {
	// Profile relates to the AWS profile you'd set.
	// In MoJ Cloud Platform this is standardised as "moj-cp".
	Profile string
	// Key refers to the AWS_ACCESS_KEY_ID, usually set in an environment variable.
	Key string
	// Secret is the AWS_SECRET_ACCESS_KEY for an AWS account.
	Secret string
	// Region relates to the AWS Region you wish to use. The default for MoJ
	// Cloud Platform should be "eu-west-2".
	Region string
	// Bucket name that holds credential file, if you wish to use one.
	Bucket string
}

// NewKubeClientWithValues takes the path of a kubeconfig file and a Kubernetes context
// as a string and returns a populated KubeClient type with the Kubernetes Clientset.
func NewKubeClientWithValues(configFilePath, context string) (*KubeClient, error) {
	// You can still build with empty strings here, but you may not get far.
	client := &KubeClient{
		Path:    configFilePath,
		Context: context,
	}

	err := client.BuildClientSet()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewMetricsKubeClientWithValues takes the path of a kubeconfig file and a Kubernetes context
// and returns a populated KubeClient type with the Kubernetes versioned.Clientset, which allows you
// to communicate with the kube metrics api.
func NewMetricsKubeClientWithValues(configFilePath, context string) (*KubeClient, error) {
	// You can still build with empty strings here, but you may not get far.
	client := &KubeClient{
		Path:    configFilePath,
		Context: context,
	}

	err := client.BuildVersionedClientset()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// BuildClientSet method builds a clientset depending on whether
// the KubeClient type has a non empty Context defined.
// If defined the clientset will have the context set to its value.
func (kube *KubeClient) BuildClientSet() (err error) {
	var config *rest.Config
	if kube.Context == "" {
		config, err = clientcmd.BuildConfigFromFlags("", kube.Path)
		if err != nil {
			return err
		}
	} else {
		config, err = kube.NewConfigFromContext()
		if err != nil {
			return err
		}
	}

	kube.Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}

// BuildClientSet method builds a versioned clientset depending on whether
// the KubeClient type has a non empty Context defined.
// If defined the clientset will have the context set to its value.
// A versioned clientset is used to communicate with the metrics api.
func (kube *KubeClient) BuildVersionedClientset() (err error) {
	var config *rest.Config
	if kube.Context == "" {
		config, err = clientcmd.BuildConfigFromFlags("", kube.Path)
		if err != nil {
			return err
		}
	} else {
		config, err = kube.NewConfigFromContext()
		if err != nil {
			return err
		}
	}

	kube.VersionedClientset, err = versioned.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}

// BuildClientSetFromS3 takes a string representing the kubeconfig file to download.
// The method uses an AwsOptions type to download said kubeconfig to the path set by the KubeClient. It
// then builds and sets a clientset.
func (kube *KubeClient) BuildClientSetFromS3(filepath string, awsOpt AwsOptions) error {
	_, err := DownloadS3Kubeconfig(filepath, kube.Path, awsOpt)
	if err != nil {
		return err
	}

	err = kube.BuildClientSet()
	if err != nil {
		return err
	}

	return nil
}

// NewConfigFromContext returns a config type ready to create a Kubernetes interface.
func (kube *KubeClient) NewConfigFromContext() (*rest.Config, error) {
	if kube.Path == "" || kube.Context == "" {
		return nil, fmt.Errorf("config file and Context must be set")
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kube.Path},
		&clientcmd.ConfigOverrides{
			CurrentContext: kube.Context,
		}).ClientConfig()
}

// SwitchKubeContext takes the KubeClient context type and ensures it
// is set appropriately.
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

// DownloadS3Kubeconfig takes a file name, the path of a kubeconfig file and an AwsOptions type.
// It will download the fileName in the awsOpts.Bucket and create a kubeconfig file in the path
// specified. The function returns the location of the newly created kubconfig.
func DownloadS3Kubeconfig(fileName, kubeconfig string, awsOpt AwsOptions) (string, error) {
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
	err = ioutil.WriteFile(kubeconfig, data, 0o644)
	if err != nil {
		return "", err
	}

	return kubeconfig, nil
}
