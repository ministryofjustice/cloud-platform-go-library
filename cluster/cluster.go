package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	VpcId      string
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

// EC2DescribeVpcEndpointConnectionsAPI defines the interface for the DescribeVpcEndpointConnections function.
// We use this interface to test the function using a mocked service.
type EC2DescribeVpcEndpointConnectionsAPI interface {
	DescribeVpcEndpointConnections(ctx context.Context,
		params *ec2.DescribeVpcEndpointConnectionsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeVpcEndpointConnectionsOutput, error)
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
	if c.Name == "live" || c.Name == "manager" {
		return errors.New("cannot create a cluster with the name live or manager")
	}

	repoName, err := findTopLevelGitDir(".")
	if err != nil {
		return err
	}

	if !strings.Contains(repoName, "cloud-platform-infrastructure") {
		return errors.New("must be run from the cloud-platform-infrastructure repository")
	}

	fmt.Printf("Creating cluster %s\n", c.Name)
	i := install.NewInstaller()

	v0_14_8 := version.Must(version.NewVersion("0.14.8"))

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
	fmt.Println("Creating VPC")
	err = c.CreateVpc(opts, execPath)
	if err != nil {
		return err
	}

	// create kubernetes cluster
	fmt.Println("Creating Kubernetes cluster")
	err = c.CreateCluster(opts, execPath)
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

func deleteLocalState(path string) error {
	if _, err := os.Stat(path); err == nil {
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func terraformInit(tf *tfexec.Terraform) error {
	err := tf.Init(context.Background())
	if err != nil {
		return err
	}
	return nil
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

func terraformPlan(tf *tfexec.Terraform, planOptions []tfexec.PlanOption, planPath, workspace string, output bool) error {
	_, err := tf.Plan(context.Background(), planOptions...)
	if err != nil {
		return err
	}

	if !output {
		return nil
	}

	plan, err := tf.ShowPlanFileRaw(context.Background(), planPath)
	if err != nil {
		return err
	}

	fmt.Println(plan)

	return nil
}

func terraformApply(tf *tfexec.Terraform, applyOptions []tfexec.ApplyOption, workspace string) error {
	var noInitErr *tfexec.ErrNoInit
	var couldNotLoad *tfexec.ErrConfigInvalid

	err := tf.Apply(context.Background(), applyOptions...)
	// handle a case where you need to init again
	if errors.As(err, &noInitErr) || errors.As(err, &couldNotLoad) {
		fmt.Println("Init required, running init again")
		err = terraformInit(tf)
		if err != nil {
			return err
		}
		err = tf.WorkspaceSelect(context.Background(), workspace)
	}
	if err != nil {
		return err
	}

	return nil
}

// GetConnectionInfo retrieves information about your VPC endpoint connections.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If successful, a DescribeVpcEndpointConnectionsOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to DescribeVpcEndpointConnections.
func GetConnectionInfo(c context.Context,
	api EC2DescribeVpcEndpointConnectionsAPI,
	input *ec2.DescribeVpcEndpointConnectionsInput) (*ec2.DescribeVpcEndpointConnectionsOutput, error) {
	return api.DescribeVpcEndpointConnections(context.Background(), input)
}

func (c *Cluster) CheckVpc(name, region string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeVpcEndpointConnectionsInput{Filters: []types.Filter{
		{Name: aws.String("vpc-endpoint-id"), Values: []string{name}},
	}}

	resp, err := GetConnectionInfo(context.Background(), client, input)
	if err != nil {
		return fmt.Errorf("error retrieving information about your VPC endpoint: %v", err)
	}

	fmt.Println("VPC endpoint connection information:", resp.VpcEndpointConnections)
	cons := len(resp.VpcEndpointConnections)

	if cons == 0 {
		return fmt.Errorf("could not find any VCP endpoint connections in " + region)
	}

	fmt.Println("VPC endpoint: Details:")
	respDecrypted, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Println(string(respDecrypted))

	fmt.Println("Found " + strconv.Itoa(cons) + " VCP endpoint connection(s) in " + region)
	return nil
}

// CreateVpc
func (c *Cluster) CreateVpc(opts *CreateOptions, execPath string) error {
	const vpcTerraformDir = "terraform/aws-accounts/cloud-platform-aws/vpc"

	fmt.Println("Checking out tf dir")
	tf, err := tfexec.NewTerraform(vpcTerraformDir, execPath)
	if err != nil {
		return fmt.Errorf("failed to create terraform: %w", err)
	}

	// if .terraform.tfstate directory exists, delete it
	fmt.Println("Deleting .terraform.tfstate directory")
	err = deleteLocalState(strings.Join([]string{vpcTerraformDir, ".terraform"}, "/"))
	if err != nil {
		return fmt.Errorf("failed to delete .terraform.tfstate directory: %w", err)
	}

	fmt.Println("Performing a terraform init")
	err = terraformInit(tf)
	if err != nil {
		return fmt.Errorf("failed to init terraform: %w", err)
	}

	fmt.Println("Creating a new workspace")
	ws, err := terraformWorkspace(c.Name, tf)
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	fmt.Println("Planning in workspace", ws)
	planPath := fmt.Sprintf("%s/%s-%v", "./", "plan", time.Now().Unix())
	planOptions := []tfexec.PlanOption{
		tfexec.Out(planPath),
		tfexec.Refresh(true),
		tfexec.Parallelism(1),
	}
	defer os.Remove(strings.Join([]string{vpcTerraformDir, planPath}, "/"))

	err = terraformPlan(tf, planOptions, planPath, ws, true)
	if err != nil {
		return fmt.Errorf("failed to plan: %w", err)
	}

	fmt.Println("Applying plan, may take a while...")
	applyOptions := []tfexec.ApplyOption{
		tfexec.DirOrPlan(planPath),
		tfexec.Parallelism(1),
	}
	err = terraformApply(tf, applyOptions, c.Name)
	if err != nil {
		return fmt.Errorf("failed to apply: %w", err)
	}

	// Get the endpoint id
	j, err := tf.Show(context.Background())
	if err != nil {
		return fmt.Errorf("failed to show: %w", err)
	}

	var vpcEndpointId string
	for k, v := range j.Values.Outputs {
		fmt.Println(k, v)
		if k == "vpc_id" {
			vpcEndpointId = v.Value.(string)
		}
	}
	if vpcEndpointId == "" {
		return fmt.Errorf("failed to find vpc endpoint id")
	}
	fmt.Println("VPC endpoint id: " + vpcEndpointId)

	// Check the vpc is created and exists
	// err = c.CheckVpc(vpcEndpointId, "eu-west-2")
	// if err != nil {
	// 	return fmt.Errorf("failed to check the vpc is up and running: %w", err)
	// }

	fmt.Println("Complete")

	return nil
}

// CreateCluster creates a new Kubernetes cluster in AWS.
func (c *Cluster) CreateCluster(opts *CreateOptions, execPath string) error {
	const eksTerraformDir = "terraform/aws-accounts/cloud-platform-aws/vpc/eks"

	fmt.Println("Checking out tf dir")
	tf, err := tfexec.NewTerraform(eksTerraformDir, execPath)
	if err != nil {
		return fmt.Errorf("failed to create terraform: %w", err)
	}

	// if .terraform.tfstate directory exists, delete it
	fmt.Println("Deleting .terraform.tfstate directory")
	err = deleteLocalState(strings.Join([]string{eksTerraformDir, ".terraform"}, "/"))
	if err != nil {
		return fmt.Errorf("failed to delete .terraform.tfstate directory: %w", err)
	}

	fmt.Println("Performing a terraform init")
	err = terraformInit(tf)
	if err != nil {
		return fmt.Errorf("failed to init terraform: %w", err)
	}

	fmt.Println("Creating a new workspace")
	ws, err := terraformWorkspace(c.Name, tf)
	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	fmt.Println("Planning in workspace", ws)
	planPath := fmt.Sprintf("%s/%s-%v", "./", "plan", time.Now().Unix())
	planOptions := []tfexec.PlanOption{
		tfexec.Out(planPath),
		tfexec.Refresh(true),
		tfexec.Parallelism(1),
	}

	err = terraformPlan(tf, planOptions, planPath, ws, true)
	if err != nil {
		return fmt.Errorf("failed to plan: %w", err)
	}

	fmt.Println("Applying plan, may take a while...")
	applyOptions := []tfexec.ApplyOption{
		tfexec.DirOrPlan(planPath),
		tfexec.Parallelism(1),
	}
	err = terraformApply(tf, applyOptions, c.Name)
	if err != nil {
		return fmt.Errorf("failed to apply: %w", err)
	}

	// Apply a tactical psp fix for the cluster
	fmt.Println("Applying a tactical psp fix for the cluster")
	err = c.ApplyTacticalPspFix()
	if err != nil {
		return fmt.Errorf("failed to apply tactical psp fix: %w", err)
	}

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
