package kubernetes

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mitchellh/mapstructure"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/nodegroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func extractKubernetesGroupMap(nodeGroups []interface{}) ([]nodegroups.NodeGroup, error) {
	filledNodeGroups := make([]nodegroups.NodeGroup, len(nodeGroups))
	for i, ng := range nodeGroups {
		g := ng.(map[string]interface{})
		for k, v := range g {
			if v == 0 {
				delete(g, k)
			}
		}
		var nodeGroup nodegroups.NodeGroup
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

func extractNodeGroupLabelsList(v []interface{}) ([]nodegroups.Label, error) {
	labels := make([]nodegroups.Label, len(v))
	for i, label := range v {
		var L nodegroups.Label
		err := mapstructure.Decode(label.(map[string]interface{}), &L)
		if err != nil {
			return nil, err
		}
		labels[i] = L
	}
	return labels, nil
}

func extractNodeGroupTaintsList(v []interface{}) ([]nodegroups.Taint, error) {
	taints := make([]nodegroups.Taint, len(v))
	for i, taint := range v {
		var T nodegroups.Taint
		err := mapstructure.Decode(taint.(map[string]interface{}), &T)
		if err != nil {
			return nil, err
		}
		taints[i] = T
	}
	return taints, nil
}

func flattenNodeGroupLabelsList(v []nodegroups.Label) []map[string]interface{} {
	labels := make([]map[string]interface{}, len(v))
	for i, label := range v {
		m := map[string]interface{}{"key": label.Key, "value": label.Value}
		labels[i] = m
	}
	return labels
}

func flattenNodeGroupTaintsList(v []nodegroups.Taint) []map[string]interface{} {
	taints := make([]map[string]interface{}, len(v))
	for i, taint := range v {
		m := map[string]interface{}{"key": taint.Key, "value": taint.Value, "effect": taint.Effect}
		taints[i] = m
	}
	return taints
}

func kubernetesStateRefreshFunc(client *gophercloud.ServiceClient, clusterID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, err := clusters.Get(client, clusterID).Extract()
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
