package cluster

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/ministryofjustice/cloud-platform-go-library/client"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Cluster struct represents an MoJ Cloud Platform Kubernetes cluster object
type Cluster struct {
	Name  string
	VpcId string

	NewestNode v1.Node
	Nodes      v1.NodeList
	OldestNode v1.Node

	Pods      *v1.PodList
	StuckPods []*v1.Pod

	Namespaces v1.NamespaceList
}

// CreateOptions struct represents the options passed to the Create method.
type CreateOptions struct {
	// Name is the name of the cluster.
	Name string
	// ClusterSuffix is the suffix to append to the cluster name.
	// This will be used to create the cluster ingress, such as "live.service.justice.gov.uk".
	ClusterSuffix string

	// NodeCount is the number of nodes to create in the cluster.
	NodeCount int
	// VpcName is the name of the VPC to create the cluster in.
	// Often clusters will be built in a single VPC.
	VpcName string

	// MaxNameLength is the maximum length of the cluster name.
	// This limit exists due to the length of the name of the ingress.
	MaxNameLength int
	// TimeOut is the maximum time to wait for the cluster to be created.
	TimeOut int
	// Debug is true if the cluster should be created in debug mode.
	Debug bool
	// Fast creates the fastest possible cluster.
	Fast bool

	// Auth0 is the Auth0 domain and secret information.
	Auth0 AuthOpts
	// AwsCredentials contains the AWS credentials to use when creating the cluster.
	AwsCredentials client.AwsCredentials

	// TerraformOptions are the options to pass to Terraform plan and apply.
	TerraformOptions TerraformOptions
	Logger           log.Logger
}

type TerraformOptions struct {
	Apply    []tfexec.ApplyOption
	Plan     []tfexec.PlanOption
	PlanPath string
	// Version is the version of Terraform to use.
	Version string
	// ExecPath is the path to the Terraform executable.
	ExecPath string
}

// AuthOpts represents the options for Auth0.
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

// verifyClusterOptions verifies the options passed to the Create method.
func verifyClusterOptions(name string, options CreateOptions) error {
	// Check the name isn't impacting a production cluster.
	if name == "live" || name == "manager" {
		return errors.New("cannot create a cluster with the name live or manager")
	}

	// Ensure the executor is running the command in the correct directory.
	repoName, err := findTopLevelGitDir(".")
	if err != nil {
		return fmt.Errorf("cannot find top level git dir: %s", err)
	}

	if !strings.Contains(repoName, "cloud-platform-infrastructure") {
		return errors.New("must be run from the cloud-platform-infrastructure repository")
	}

	return nil
}

// createTerraformObj creates a Terraform object using the version passed as a string.
func (tf *TerraformOptions) CreateTerraformObj() error {
	i := install.NewInstaller()
	v := version.Must(version.NewVersion(tf.Version))

	execPath, err := i.Ensure(context.Background(), []src.Source{
		&fs.ExactVersion{
			Product: product.Terraform,
			Version: v,
		},
		&releases.ExactVersion{
			Product: product.Terraform,
			Version: v,
		},
	})
	if err != nil {
		return err
	}

	defer i.Remove(context.Background())

	tf.ExecPath = execPath
	// Define the Terraform options.
	tf.PlanPath = fmt.Sprintf("%s/%s-%v", "./", "plan", time.Now().Unix())
	tf.Plan = []tfexec.PlanOption{
		tfexec.Out(tf.PlanPath),
		tfexec.Refresh(true),
		tfexec.Parallelism(1),
	}
	tf.Apply = []tfexec.ApplyOption{
		tfexec.DirOrPlan(tf.PlanPath),
		tfexec.Parallelism(1),
	}

	return nil
}

// Create creates a new Kubernetes cluster using the options passed to it.
func (c *Cluster) Create(opts *CreateOptions) error {
	// Directory paths in the cloud-platform-infrastructure repository.
	const (
		baseDir       = "terraform/aws-accounts/cloud-platform-aws/"
		vpcDir        = baseDir + "vpc/"
		clusterDir    = vpcDir + "eks/"
		componentsDir = clusterDir + "components/"
	)

	// Verify the options passed to Create are allowed.
	err := verifyClusterOptions(opts.Name, *opts)
	if err != nil {
		return fmt.Errorf("error verifying cluster options: %s", err)
	}

	// Add name to the cluster object.
	c.Name = opts.Name

	// Create Terraform object to use throught method.
	err = opts.TerraformOptions.CreateTerraformObj()
	if err != nil {
		return fmt.Errorf("error creating terraform obj: %s", err)
	}

	// Create VPC
	state, err := c.TerraformApply(opts, vpcDir)
	if err != nil {
		return fmt.Errorf("error creating vpc: %s", err)
	}

	// Check the vpc is created and exists
	err = c.CheckVpc(state, opts.AwsCredentials.Session)
	if err != nil {
		return fmt.Errorf("failed to check the vpc is up and running: %w", err)
	}

	// If the user specifies a fast build, then we don't need to create the auth0 module.
	if opts.Fast {
		opts.TerraformOptions.Apply = append(opts.TerraformOptions.Apply, tfexec.Var(fmt.Sprintf("%s=%v", "auth0_count", 0)))
		fmt.Println("Fast mode enabled, skipping auth0 creation")
	}

	// Create the Kubernetes cluster.
	fmt.Println("Creating Kubernetes cluster")
	// _, err = c.TerraformApply(opts, clusterDir)
	// if err != nil {
	// 	return err
	// }

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

func getVpcFromState(state *tfjson.State) (string, error) {
	var vpcEndpointId string
	for k, v := range state.Values.Outputs {
		if k == "vpc_id" {
			vpcEndpointId = v.Value.(string)
		}
	}
	if vpcEndpointId == "" {
		return "", fmt.Errorf("failed to find vpc endpoint id")
	}

	return vpcEndpointId, nil
}

func (c *Cluster) terraformInitApply(dir string, tf *tfexec.Terraform, opts CreateOptions) (*tfjson.State, error) {
	ws, err := intialise(c.Name, tf)
	if err != nil {
		return nil, fmt.Errorf("failed to init terraform: %w", err)
	}

	fmt.Println("Planning in workspace", ws)
	defer os.Remove(strings.Join([]string{dir, opts.TerraformOptions.PlanPath}, "/"))

	err = plan(tf, opts.TerraformOptions, true)
	if err != nil {
		return nil, fmt.Errorf("failed to plan: %w", err)
	}

	fmt.Println("Applying plan, may take a while...")
	applyOptions := []tfexec.ApplyOption{}
	err = apply(tf, applyOptions, c.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to apply: %w", err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to show: %w", err)
	}
	fmt.Println("Vpc Complete")

	return state, nil
}

func deleteLocalState(path string) error {
	if _, err := os.Stat(path); err == nil {
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	return nil
}

// intialise performs the `terraform init` function.
func intialise(workspace string, tf *tfexec.Terraform) (string, error) {
	err := tf.Init(context.Background())
	if err != nil {
		return "", err
	}
	return terraformWorkspace(workspace, tf)
}

func terraformWorkspace(workspace string, tf *tfexec.Terraform) (string, error) {
	list, _, err := tf.WorkspaceList(context.Background())

	for _, ws := range list {
		if ws == workspace {
			err = tf.WorkspaceSelect(context.Background(), workspace)
			if err != nil {
				return "", err
			}
			return workspace, nil
		}
	}

	err = tf.WorkspaceNew(context.Background(), workspace)
	if err != nil {
		return "", err
	}
	ws, err := tf.WorkspaceShow(context.Background())
	if err != nil {
		return "", err
	}

	return ws, nil
}

func plan(tf *tfexec.Terraform, opt TerraformOptions, output bool) error {
	_, err := tf.Plan(context.Background(), opt.Plan...)
	if err != nil {
		return fmt.Errorf("failed to execute the plan command: %w", err)
	}

	if !output {
		return nil
	}

	plan, err := tf.ShowPlanFileRaw(context.Background(), opt.PlanPath)
	if err != nil {
		return fmt.Errorf("failed to show the plan file: %w", err)
	}

	fmt.Println(plan)

	return nil
}

func apply(tf *tfexec.Terraform, applyOptions []tfexec.ApplyOption, workspace string) error {
	var noInitErr *tfexec.ErrNoInit
	var couldNotLoad *tfexec.ErrConfigInvalid

	err := tf.Apply(context.Background(), applyOptions...)
	// handle a case where you need to init again
	if errors.As(err, &noInitErr) || errors.As(err, &couldNotLoad) {
		fmt.Println("Init required, running init again")
		_, err = intialise(workspace, tf)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	return nil
}

// CheckVpc asserts that the vpc is up and running. It tests the vpc state and id.
func (c *Cluster) CheckVpc(state *tfjson.State, sess *session.Session) error {
	vpcId, err := getVpcFromState(state)
	if err != nil {
		return fmt.Errorf("unable to get vpcid from statefile: %e", err)
	}

	svc := ec2.New(sess)

	vpc, err := svc.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Cluster"),
				Values: []*string{aws.String(c.Name)},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error describing vpc: %v", err)
	}

	if len(vpc.Vpcs) == 0 {
		return fmt.Errorf("no vpc found")
	}

	if vpc.Vpcs[0].VpcId != nil && *vpc.Vpcs[0].VpcId != vpcId {
		return fmt.Errorf("vpc id mismatch: %s != %s", *vpc.Vpcs[0].VpcId, vpcId)
	}

	if vpc.Vpcs[0].State != nil && *vpc.Vpcs[0].State != "available" {
		return fmt.Errorf("vpc not available: %s", *vpc.Vpcs[0].State)
	}

	c.VpcId = *vpc.Vpcs[0].VpcId

	return nil
}

// CreateVpc
func (c *Cluster) TerraformApply(opts *CreateOptions, directory string) (*tfjson.State, error) {
	var out bytes.Buffer

	fmt.Println("Checking out tf dir")
	tf, err := tfexec.NewTerraform(directory, opts.TerraformOptions.ExecPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create terraform: %w", err)
	}

	tf.SetStdout(&out)
	tf.SetStderr(&out)

	// if .terraform.tfstate directory exists, delete it
	fmt.Println("Deleting .terraform.tfstate directory")
	err = deleteLocalState(strings.Join([]string{directory, ".terraform"}, "/"))
	if err != nil {
		return nil, fmt.Errorf("failed to delete .terraform.tfstate directory: %w", err)
	}

	return c.terraformInitApply(directory, tf, *opts)
}

// InstallComponents installs components into the Kubernetes cluster.
func installComponents(opts *CreateOptions) error {
	return nil
}

// HealthCheck performs a health check on the Kubernetes cluster.
func healthCheck(opts *CreateOptions) error {
	return nil
}

// ApplyTacticalPspFix deletes the current eks.privileged psp in the cluster.
// This allows the cluster to be created with a different psp. All pods are recycled
// so the new psp will be applied.
func (c *Cluster) ApplyTacticalPspFix() error {
	client, err := client.NewKubeClientWithValues("", "")
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Delete the eks.privileged psp
	err = client.Clientset.PolicyV1beta1().PodSecurityPolicies().Delete(context.Background(), "eks.privileged", metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete eks.privileged psp: %w", err)
	}

	// Delete all pods in the cluster
	err = client.Clientset.CoreV1().Pods("").DeleteCollection(context.Background(), metav1.DeleteOptions{}, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to recycle pods: %w", err)
	}

	return nil
}
