package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
)

func TestCreateOptsAutoScaleZeroSize(t *testing.T) {
	opts := nodegroups.CreateOpts{
		ClusterID: "cluster-id",
		Spec: nodegroups.NodeGroupSpec{
			Name: "auto",
			VMEngine: nodegroups.VMEngine{
				NovaEngine: nodegroups.NovaEngine{FlavorID: "flavor-id"},
			},
			Zones: []string{"MS1"},
			ScaleSpec: nodegroups.ScaleSpec{
				AutoScale: &nodegroups.AutoScale{MinSize: 0, MaxSize: 5, Size: 0},
			},
			DiskType: nodegroups.DiskType{
				CinderVolumeType: nodegroups.CinderVolumeType{Type: "ceph-ssd", Size: 128},
			},
		},
	}

	b, err := opts.ToNodeGroupCreateMap()
	require.NoError(t, err)

	spec := b["node_group_spec"].(map[string]interface{})
	autoScale := spec["scale_spec"].(map[string]interface{})["auto_scale"].(map[string]interface{})
	assert.Equal(t, map[string]interface{}{"min_size": float64(0), "max_size": float64(5), "size": float64(0)}, autoScale)
}

func TestUpdateOptsAutoScaleZeroSize(t *testing.T) {
	opts := nodegroups.UpdateOpts{
		ScaleSpec: &nodegroups.ScaleSpec{
			AutoScale: &nodegroups.AutoScale{MinSize: 0, MaxSize: 3, Size: 0},
		},
	}

	b, err := opts.ToNodeGroupUpdateMap()
	require.NoError(t, err)

	autoScale := b["scale_spec"].(map[string]interface{})["auto_scale"].(map[string]interface{})
	assert.Equal(t, map[string]interface{}{"min_size": float64(0), "max_size": float64(3), "size": float64(0)}, autoScale)
}
