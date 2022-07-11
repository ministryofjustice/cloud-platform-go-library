package namespace

import (
	"context"
	"fmt"

	"github.com/ministryofjustice/cloud-platform-go-library/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AllNamespaces returns all namespaces in the cluster.
func AllNamespaces(c *client.KubeClient) (*v1.NamespaceList, error) {
	list, err := c.Clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Namespace returns the namespace with the given name.
func Namespace(c *client.KubeClient, name string) (*v1.Namespace, error) {
	// Gather list of namespaces from a cluster
	list, err := AllNamespaces(c)
	if err != nil {
		return nil, err
	}
	for _, namespace := range list.Items {
		if namespace.Name == name {
			return &namespace, nil
		}
	}

	return nil, fmt.Errorf("namespace %s not found", name)
}

// CreateNamespace creates a new namespace with the given name.
func CreateNamespace(c *client.KubeClient, name string) (*v1.Namespace, error) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return c.Clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
}

// DeleteNamespace deletes the namespace with the given name.
func DeleteNamespace(c *client.KubeClient, name string) error {
	ns, err := Namespace(c, name)
	if err != nil {
		return err
	}
	return c.Clientset.CoreV1().Namespaces().Delete(context.Background(), ns.Name, metav1.DeleteOptions{})
}

// GetTeamNamespaces returns all namespaces in the cluster that are owned by the given team.
// A team is a GitHub group or organization defined in the cloud-platform environments repository.
func GetTeamNamespaces(team string) ([]*v1.Namespace, error) {
	return nil, nil
}

// NamespaceSlackChannel returns the slack channel name for the given namespace.
func NamespaceSlackChannel(c *client.KubeClient, name string) (string, error) {
	ns, err := Namespace(c, name)
	if err != nil {
		return "", err
	}
	return ns.Annotations["cloud-platform.justice.gov.uk/slack-channel"], nil
}

// ProductionNamespace returns a slice of namespaces with a production label.
func ProductionNamespace(c *client.KubeClient) ([]*v1.Namespace, error) {
	list, err := AllNamespaces(c)
	if err != nil {
		return nil, err
	}
	var namespaces []*v1.Namespace
	for _, namespace := range list.Items {
		if namespace.Labels["cloud-platform.justice.gov.uk/is-production"] == "true" {
			namespaces = append(namespaces, &namespace)
		}
	}
	return namespaces, nil
}

// NonProductionNamespace returns a slice of namespaces without a production label.
func NonProductionNamespace(c *client.KubeClient) ([]*v1.Namespace, error) {
	list, err := AllNamespaces(c)
	if err != nil {
		return nil, err
	}
	var namespaces []*v1.Namespace
	for _, namespace := range list.Items {
		if namespace.Labels["cloud-platform.justice.gov.uk/is-production"] != "true" {
			namespaces = append(namespaces, &namespace)
		}
	}
	return namespaces, nil
}

// NamespaceGithubTeam returns the github team name for the given namespace.
func NamespaceSourceCode(c *client.KubeClient, name string) (string, error) {
	ns, err := Namespace(c, name)
	if err != nil {
		return "", err
	}
	return ns.Annotations["cloud-platform.justice.gov.uk/source-code"], nil
}

// NamespaceOwner returns the owner of the namespace.
func NamespaceOwner(c *client.KubeClient, name string) (string, error) {
	ns, err := Namespace(c, name)
	if err != nil {
		return "", err
	}
	return ns.Annotations["cloud-platform.justice.gov.uk/team-name"], nil
}
