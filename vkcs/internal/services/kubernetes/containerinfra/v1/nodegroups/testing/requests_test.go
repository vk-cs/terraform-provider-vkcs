package testing

import (
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v1/nodegroups"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatchOpts(t *testing.T) {

	patchOpts := nodegroups.PatchOpts{
		{
			Path:  "/max_nodes",
			Value: 10,
			Op:    "replace",
		},
		{
			Path:  "/min_nodes",
			Value: 2,
			Op:    "replace",
		},
		{
			Path:  "/autoscaling_enabled",
			Value: false,
			Op:    "replace",
		},
		{
			Path:  "/max_node_unavailable",
			Value: 80,
			Op:    "replace",
		},
	}

	b, _ := patchOpts.PatchMap()

	assert.IsType(t, []map[string]interface{}{}, b)
	assert.Len(t, b, 4)
}

func TestAddBatchOpts(t *testing.T) {

	addGroups := []nodegroups.NodeGroup{
		{
			Name:     "test1",
			FlavorID: "95663bae-6763-4a53-9424-831975285cc1",
		},
		{
			Name:     "test2",
			FlavorID: "95663bae-6763-4a53-9424-831975285cc1",
		},
	}

	addbatchOpts := nodegroups.BatchAddParams{
		Action:  "batch_add_ng",
		Payload: addGroups,
	}

	b, _ := addbatchOpts.Map()

	assert.IsType(t, []interface{}{}, b["payload"])
	assert.Len(t, b["payload"], 2)
	assert.Len(t, b, 2)
}
