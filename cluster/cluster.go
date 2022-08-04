package cluster

import (
	"github.com/ministryofjustice/cloud-platform-go-library/client"
	v1 "k8s.io/api/core/v1"
)

// Cluster struct represents an MoJ Cloud Platform Kubernetes cluster object
type Cluster struct {
	Name       string
	NewestNode v1.Node
	Nodes      v1.NodeList
	OldestNode v1.Node
	Pods       *v1.PodList
	StuckPods  []*v1.Pod
	Namespaces v1.NamespaceList
}

// CreateOptions struct represents the options passed to the Create method.
type CreateOptions struct {
	Name          string
	ClusterSuffix string

	NodeCount int
	VpcName   string

	Kubeconfig string

	MaxNameLength int `default:"12"`
	TimeOut       int
	Debug         bool

	Auth0          AuthOpts
	AwsCredentials client.AwsCredentials
}

type AuthOpts struct {
	Domain       string
	ClientId     string
	ClientSecret string
}

// NewWithValues returns a full Cluster object with populated values.
func NewWithValues(c client.KubeClient) (*Cluster, error) {
	nodes, err := AllNodes(c)
	if err != nil {
		return nil, err
	}
	pods, err := AllPods(c)
	if err != nil {
		return nil, err
	}
	stuckPods, err := StuckPods(c, *pods)
	if err != nil {
		return nil, err
	}
	cluster := &Cluster{
		Nodes:     nodes,
		Pods:      pods,
		StuckPods: stuckPods,
	}

	// You can only get the name of a Cloud Platform cluster using the labels on a node.
	cluster.GetName()

	return cluster, nil
}

// GetName is a method function to get the name of the cluster.
func (c *Cluster) GetName() {
	c.Name = c.Nodes.Items[0].Labels["Cluster"]
}

// Create creates a new Kubernetes cluster using the options passed to it.
func (c *Cluster) Create(opts *CreateOptions) error {
	// create vpc
	err := createVpc(opts)
	if err != nil {
		return err
	}

	// create kubernetes cluster
	err = createCluster(opts)
	if err != nil {
		return err
	}

	// install components into kubernetes cluster
	err = installComponents(opts)
	if err != nil {
		return err
	}

	// perform health check on the cluster
	err = healthCheck(opts)
	if err != nil {
		return err
	}

	return nil
}

// createVpc creates a new VPC in AWS.
func createVpc(opts *CreateOptions) error {
	return nil
}

// CreateCluster creates a new Kubernetes cluster in AWS.
func createCluster(opts *CreateOptions) error {
	return nil
}

// InstallComponents installs components into the Kubernetes cluster.
func installComponents(opts *CreateOptions) error {
	return nil
}

// HealthCheck performs a health check on the Kubernetes cluster.
func healthCheck(opts *CreateOptions) error {
	return nil
}
