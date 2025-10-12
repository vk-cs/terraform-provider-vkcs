package kubernetes

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v1/nodegroups"
)

func TestExpandContainerInfraLabelsMap(t *testing.T) {
	labels := map[string]interface{}{
		"foo": "bar",
		"bar": "baz",
	}

	expectedLabels := map[string]string{
		"foo": "bar",
		"bar": "baz",
	}

	actualLabels, err := extractKubernetesLabelsMap(labels)
	assert.Equal(t, err, nil)
	assert.Equal(t, expectedLabels, actualLabels)
}

func TestExpandKubernetesGroupMap(t *testing.T) {
	ncount, maxn, minn, vs := 2, 3, 1, 10

	var groups []interface{}
	group := map[string]interface{}{
		"name":                "test",
		"node_count":          ncount,
		"max_nodes":           maxn,
		"min_nodes":           minn,
		"volume_size":         vs,
		"volume_type":         "",
		"flavor_id":           "1",
		"image_id":            "1",
		"autoscaling_enabled": false,
	}
	groups = append(groups, group)

	var expectedGroups []nodegroups.NodeGroup
	expgroup := nodegroups.NodeGroup{
		Name:       "test",
		NodeCount:  ncount,
		MaxNodes:   maxn,
		MinNodes:   minn,
		VolumeSize: vs,
		FlavorID:   "1",
		ImageID:    "1",
	}
	expectedGroups = append(expectedGroups, expgroup)

	actualGroups, err := extractKubernetesGroupMap(groups)

	assert.Equal(t, err, nil)
	assert.Equal(t, expectedGroups, actualGroups)
}

func TestExtractNodeGroupLabelsList(t *testing.T) {
	labels := []interface{}{
		map[string]interface{}{
			"key":   "bar",
			"value": "baz",
		},
		map[string]interface{}{
			"key":   "foo",
			"value": "bar",
		},
		map[string]interface{}{
			"key": "label_without_value",
		},
	}
	expectedLabels := []nodegroups.Label{
		{
			Key:   "bar",
			Value: "baz",
		},
		{
			Key:   "foo",
			Value: "bar",
		},
		{
			Key: "label_without_value",
		},
	}

	actualLabels, err := extractNodeGroupLabelsList(labels)
	sort.Slice(actualLabels, func(i, j int) bool {
		return actualLabels[i].Key < actualLabels[j].Key
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, expectedLabels, actualLabels)
}

func TestExtractNodeGroupTaintsList(t *testing.T) {
	taints := []any{
		map[string]any{
			"key":    "key1",
			"value":  "val1",
			"effect": "effect1",
		},
		map[string]any{
			"key":    "key2",
			"value":  "",
			"effect": "effect2",
		},
	}

	expectedTaints := []nodegroups.Taint{
		{
			Key:    "key1",
			Value:  "val1",
			Effect: "effect1",
		},
		{
			Key:    "key2",
			Value:  "",
			Effect: "effect2",
		},
	}

	actualTaints, err := extractNodeGroupTaintsList(taints)
	assert.Nil(t, err)
	sort.Slice(actualTaints, func(i, j int) bool {
		return actualTaints[i].Key < actualTaints[j].Key
	})
	assert.Equal(t, expectedTaints, actualTaints)
}

func TestExtractNodeGroupTaintsListError(t *testing.T) {
	taints := []any{
		nil,
		map[string]any{
			"key":    "key2",
			"value":  "",
			"effect": "effect2",
		},
	}

	actualTaints, err := extractNodeGroupTaintsList(taints)
	assert.Empty(t, actualTaints)
	assert.EqualError(t, err, "empty node group taint with index: 0")
}
