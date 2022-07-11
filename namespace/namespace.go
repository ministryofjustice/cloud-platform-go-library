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
	// Gather list of namespaces from a cluster
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
