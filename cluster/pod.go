package cluster

import (
	"context"
	"fmt"

	"github.com/ministryofjustice/cloud-platform-go-library/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AllPods returns all pods in a cluster as PodList objects.
func AllPods(c client.KubeClient) (*v1.PodList, error) {
	podList, err := c.Clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %s", err)
	}

	return podList, nil
}

// StuckPods returns all pods in a cluster that are in a non-running state.
func StuckPods(c client.KubeClient, pods v1.PodList) ([]*v1.Pod, error) {
	var stuckPods []*v1.Pod
	for _, pod := range pods.Items {
		if pod.Status.Phase == v1.PodPending || pod.Status.Phase == v1.PodFailed || pod.Status.Phase == v1.PodUnknown {
			stuckPods = append(stuckPods, &pod)
		}
	}

	return stuckPods, nil
}
