package cluster_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCluster_GetName(t *testing.T) {
	standard.Cluster.GetName()

	assert.Equal(t, "Cluster1", standard.Cluster.Name)
}
