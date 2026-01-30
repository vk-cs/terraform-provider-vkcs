package kubernetes

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	v1clusters "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v1/clusters"
	v1nodegroups "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v1/nodegroups"
	v2clusters "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	v2nodegroups "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type Schema interface {
	Get(key string) interface{}
}

// getAsTaintsSlice возвращает значение поля как []nodegroups.Taint, независимо от его типа ([]interface{} или *schema.Set)
func getAsTaintsSlice(d Schema, key string) ([]v2nodegroups.Taint, error) {
	value := d.Get(key)

	switch v := value.(type) {
	case []interface{}:
		return extractNodeGroupTaintsListV2(v)
	case *schema.Set:
		list := v.List()
		return extractNodeGroupTaintsListV2(list)
	default:
		return nil, fmt.Errorf("field %s has unsupported type %T", key, v)
	}
}

// getAsStringSlice возвращает значение поля как []string, независимо от его типа ([]interface{} или *schema.Set)
func getAsStringSlice(d Schema, key string) ([]string, error) {
	value := d.Get(key)

	switch v := value.(type) {
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				result[i] = str
			} else {
				return nil, fmt.Errorf("item at index %d in field %s is not a string", i, key)
			}
		}
		return result, nil
	case *schema.Set:
		list := v.List()
		result := make([]string, len(list))
		for i, item := range list {
			if str, ok := item.(string); ok {
				result[i] = str
			} else {
				return nil, fmt.Errorf("item at index %d in field %s is not a string", i, key)
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("field %s has unsupported type %T", key, v)
	}
}

func extractKubernetesGroupMap(nodeGroups []interface{}) ([]v1nodegroups.NodeGroup, error) {
	filledNodeGroups := make([]v1nodegroups.NodeGroup, len(nodeGroups))
	for i, ng := range nodeGroups {
		g := ng.(map[string]interface{})
		for k, v := range g {
			if v == 0 {
				delete(g, k)
			}
		}
		var nodeGroup v1nodegroups.NodeGroup
		err := util.MapStructureDecoder(&nodeGroup, &g, util.DecoderConfig)
		if err != nil {
			return nil, err
		}
		filledNodeGroups[i] = nodeGroup
	}
	return filledNodeGroups, nil
}

func extractKubernetesLabelsMap(v map[string]interface{}) (map[string]string, error) {
	m := make(map[string]string)
	for key, val := range v {
		labelValue, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("label %s value should be string", key)
		}
		m[key] = labelValue
	}
	return m, nil
}

func extractNodeGroupLabelsList(v []interface{}) ([]v1nodegroups.Label, error) {
	labels := make([]v1nodegroups.Label, len(v))
	for i, label := range v {
		var L v1nodegroups.Label
		err := mapstructure.Decode(label.(map[string]interface{}), &L)
		if err != nil {
			return nil, err
		}
		labels[i] = L
	}
	return labels, nil
}

func extractNodeGroupTaintsListV1(rawTaints []any) ([]v1nodegroups.Taint, error) {
	taints := make([]v1nodegroups.Taint, len(rawTaints))
	for i, rawTaint := range rawTaints {
		taint, ok := rawTaint.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("empty node group taint with index: %d", i)
		}

		var resTaint v1nodegroups.Taint
		if err := mapstructure.Decode(taint, &resTaint); err != nil {
			return nil, fmt.Errorf("failed to read node group taint with index %d: %s", i, err)
		}

		taints[i] = resTaint
	}

	return taints, nil
}

func extractNodeGroupTaintsListV2(rawTaints []any) ([]v2nodegroups.Taint, error) {
	taints := make([]v2nodegroups.Taint, len(rawTaints))
	for i, rawTaint := range rawTaints {
		taint, ok := rawTaint.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("empty node group taint with index: %d", i)
		}

		var resTaint v2nodegroups.Taint
		if err := mapstructure.Decode(taint, &resTaint); err != nil {
			return nil, fmt.Errorf("failed to read node group taint with index %d: %s", i, err)
		}

		taints[i] = resTaint
	}

	return taints, nil
}

func flattenNodeGroupLabelsList(v []v1nodegroups.Label) []map[string]interface{} {
	labels := make([]map[string]interface{}, len(v))
	for i, label := range v {
		m := map[string]interface{}{"key": label.Key, "value": label.Value}
		labels[i] = m
	}
	return labels
}

func flattenNodeGroupTaintsList(v []v1nodegroups.Taint) []map[string]interface{} {
	taints := make([]map[string]interface{}, len(v))
	for i, taint := range v {
		m := map[string]interface{}{"key": taint.Key, "value": taint.Value, "effect": taint.Effect}
		taints[i] = m
	}
	return taints
}

func kubernetesStateRefreshFunc(client *gophercloud.ServiceClient, clusterID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, err := v1clusters.Get(client, clusterID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return c, string(clusterStatusNotFound), nil
			}
			return nil, "", err
		}
		if c.NewStatus == string(clusterStatusError) {
			err = fmt.Errorf("vkcs_kubernetes_cluster is in an error state: %s", c.StatusReason)
			return c, c.NewStatus, err
		}
		return c, c.NewStatus, nil
	}
}

func kubernetesStateRefreshFuncV2(client *gophercloud.ServiceClient, clusterID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, err := v2clusters.Get(client, clusterID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return c, string(clusterStatusV2Deleted), nil
			}
			return nil, "", err
		}
		if c.Status == clusterStatusV2Failed {
			err = fmt.Errorf("vkcs_kubernetes_cluster_v2 is in an error state")
			return c, c.Status, err
		}
		return c, c.Status, nil
	}
}

func kubernetesNodeGroupStateRefreshFunc(client *gophercloud.ServiceClient, nodeGroupID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, err := v1nodegroups.Get(client, nodeGroupID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return c, string(nodeGroupStatusNotFound), nil
			}
			return nil, "", err
		}
		if c.State == string(nodeGroupStatusError) {
			err = fmt.Errorf("vkcs_kubernetes_node_group is in an error state")
			return c, c.State, err
		}
		return c, c.State, nil
	}
}
