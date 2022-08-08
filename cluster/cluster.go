package cluster

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/ministryofjustice/cloud-platform-go-library/client"
	"github.com/pkg/errors"
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

	MaxNameLength int
	TimeOut       int
	Debug         bool

	Auth0          AuthOpts
	AwsCredentials client.AwsCredentials

	Logger log.Logger
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

func findTopLevelGitDir(workingDir string) (string, error) {
	dir, err := filepath.Abs(workingDir)
	if err != nil {
		return "", errors.Wrap(err, "invalid working dir")
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("no git repository found")
		}
		dir = parent
	}
}

// Create creates a new Kubernetes cluster using the options passed to it.
func (c *Cluster) Create(opts *CreateOptions) error {
	c.Name = opts.Name

	repoName, err := findTopLevelGitDir(".")
	if err != nil {
		return err
	}

	if !strings.Contains(repoName, "cloud-platform-infrastructure") {
		return errors.New("must be run from the cloud-platform-infrastructure repository")
	}

	fmt.Printf("Creating cluster %s", c.Name)
	i := install.NewInstaller()

	v0_14_8 := version.Must(version.NewVersion("0.14.8"))

	fmt.Println("Installing terraform")
	execPath, err := i.Ensure(context.Background(), []src.Source{
		&fs.ExactVersion{
			Product: product.Terraform,
			Version: v0_14_8,
		},
		&releases.ExactVersion{
			Product: product.Terraform,
			Version: v0_14_8,
		},
	})

	defer i.Remove(context.Background())

	// create vpc
	err = c.CreateVpc(opts, execPath)
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
func (c *Cluster) CreateVpc(opts *CreateOptions, execPath string) error {
	fmt.Println("Checking out tf dir")
	workingDir := "terraform/aws-accounts/cloud-platform-aws/vpc"
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	// if .terraform.tfstate directory exists, delete it
	fmt.Println("Deleting .terraform.tfstate directory")
	if _, err := os.Stat("terraform/aws-accounts/cloud-platform-aws/vpc/.terraform"); err == nil {
		err = os.RemoveAll("terraform/aws-accounts/cloud-platform-aws/vpc/.terraform")
		if err != nil {
			return err
		}
	}

	fmt.Println("Performing a terraform init")
	err = tf.Init(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Creating a new workspace")
	err = tf.WorkspaceNew(context.Background(), c.Name)
	if err != nil {
		return err
	}

	ws, err := tf.WorkspaceShow(context.Background())
	fmt.Println("Applying in workspace", ws)
	planPath := fmt.Sprintf("%s/%s-%v", "./", "plan", time.Now().Unix())
	planOptions := []tfexec.PlanOption{
		tfexec.Out(planPath),
		tfexec.Refresh(true),
		tfexec.Parallelism(1),
	}
	_, err = tf.Plan(context.Background(), planOptions...)
	if err != nil {
		return err
	}

	fmt.Println("Generating plan file")
	plan, err := tf.ShowPlanFile(context.Background(), planPath)
	if err != nil {
		return err
	}

	// TODO: Make this a debug option
	for k, v := range plan.OutputChanges {
		fmt.Printf("%s: %s\n", k, v)
	}

	fmt.Println("Applying plan")
	statePath := fmt.Sprintf("%s/%s-%v", "./", "state", time.Now().Unix())
	applyOptions := []tfexec.ApplyOption{
		tfexec.DirOrPlan(planPath),
		tfexec.StateOut(statePath),
		tfexec.Parallelism(1),
	}
	err = tf.Apply(context.Background(), applyOptions...)
	if err != nil {
		return err
	}

	// Show the final terraform
	fmt.Println("Showing final state")
	state, err := tf.ShowStateFile(context.Background(), statePath)
	if err != nil {
		return err
	}

	for k, v := range state.Values.Outputs {
		fmt.Printf("%s: %v\n", k, v)
	}

	destroyOptions := []tfexec.DestroyOption{
		tfexec.Refresh(true),
		tfexec.Parallelism(1),
	}

	err = tf.Destroy(context.Background(), destroyOptions...)
	if err != nil {
		return err
	}

	err = tf.WorkspaceSelect(context.Background(), "default")
	if err != nil {
		return err
	}

	err = tf.WorkspaceDelete(context.Background(), c.Name)
	if err != nil {
		return err
	}

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
