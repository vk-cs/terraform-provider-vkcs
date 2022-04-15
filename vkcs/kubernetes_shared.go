package vkcs

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/mitchellh/mapstructure"
)

func extractKubernetesGroupMap(nodeGroups []interface{}) ([]nodeGroup, error) {
	filledNodeGroups := make([]nodeGroup, len(nodeGroups))
	for i, ng := range nodeGroups {
		g := ng.(map[string]interface{})
		for k, v := range g {
			if v == 0 {
				delete(g, k)
			}
		}
		var nodeGroup nodeGroup
		err := mapStructureDecoder(&nodeGroup, &g, decoderConfig)
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

func extractNodeGroupLabelsList(v []interface{}) ([]nodeGroupLabel, error) {
	labels := make([]nodeGroupLabel, len(v))
	for i, label := range v {
		var L nodeGroupLabel
		err := mapstructure.Decode(label.(map[string]interface{}), &L)
		if err != nil {
			return nil, err
		}
		labels[i] = L
	}
	return labels, nil
}

func extractNodeGroupTaintsList(v []interface{}) ([]nodeGroupTaint, error) {
	taints := make([]nodeGroupTaint, len(v))
	for i, taint := range v {
		var T nodeGroupTaint
		err := mapstructure.Decode(taint.(map[string]interface{}), &T)
		if err != nil {
			return nil, err
		}
		taints[i] = T
	}
	return taints, nil
}

func flattenNodeGroupLabelsList(v []nodeGroupLabel) []map[string]interface{} {
	labels := make([]map[string]interface{}, len(v))
	for i, label := range v {
		m := map[string]interface{}{"key": label.Key, "value": label.Value}
		labels[i] = m
	}
	return labels
}

func flattenNodeGroupTaintsList(v []nodeGroupTaint) []map[string]interface{} {
	taints := make([]map[string]interface{}, len(v))
	for i, taint := range v {
		m := map[string]interface{}{"key": taint.Key, "value": taint.Value, "effect": taint.Effect}
		taints[i] = m
	}
	return taints
}

func kubernetesStateRefreshFunc(client ContainerClient, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, err := clusterGet(client, clusterID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return c, string(clusterStatusDeleted), nil
			}
			return nil, "", err
		}
		if c.NewStatus == clusterStatusError {
			err = fmt.Errorf("vkcs_kubernetes_cluster is in an error state: %s", c.StatusReason)
			return c, string(c.NewStatus), err
		}
		return c, string(c.NewStatus), nil
	}
}
