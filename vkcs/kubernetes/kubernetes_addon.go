package kubernetes

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	v1 "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfraaddons/v1"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfraaddons/v1/addons"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfraaddons/v1/clusteraddons"
)

func readAllClusterAddons(ctx context.Context, client *gophercloud.ServiceClient, clusterID string) ([]v1.Addon, []clusteraddons.ClusterAddon, error) {
	tflog.Debug(ctx, "Calling Addons API to list available addons")

	allAvailableAddons, err := addons.ListClusterAvailableAddons(client, clusterID).AllPages()
	if err != nil {
		return nil, nil, fmt.Errorf("error listing addons available for the cluster: %s", err)
	}

	availableAddons, err := addons.ExtractAddons(allAvailableAddons)
	if err != nil {
		return nil, nil, fmt.Errorf("error extracting addons available for the cluster: %w", err)
	}

	tflog.Debug(ctx, "Called Addons API to list available addons", map[string]interface{}{"available_addons": fmt.Sprintf("%#v", availableAddons)})
	tflog.Debug(ctx, "Calling Addons API to list cluster addons")

	allClusterAddons, err := addons.ListClusterAddons(client, clusterID).AllPages()
	if err != nil {
		return availableAddons, nil, fmt.Errorf("error listing addons installed in the cluster: %w", err)
	}

	clusterAddons, err := addons.ExtractClusterAddons(allClusterAddons)
	if err != nil {
		return availableAddons, nil, fmt.Errorf("error extracting addons installed in the cluster: %w", err)
	}

	tflog.Debug(ctx, "Called Addons API to list cluster addons", map[string]interface{}{"cluster_addons": fmt.Sprintf("%#v", clusterAddons)})

	return availableAddons, clusterAddons, nil
}

func addonStateRefreshFunc(client *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		a, err := clusteraddons.Get(client, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return a, addonStatusDeleted, nil
			}
			if a.Status == addonStatusFailed {
				return a, a.Status, fmt.Errorf("there was an error creating the kubernetes cluster addon: %s", err)
			}
			return nil, "", err
		}
		return a, a.Status, nil
	}
}

func addonIn(addon v1.Addon, addons []v1.Addon) bool {
	for _, a := range addons {
		if a.ID == addon.ID {
			return true
		}
	}
	return false
}

func extractAddonsFromClusterAddons(clusterAddons []clusteraddons.ClusterAddon) ([]v1.Addon, error) {
	addons := make([]v1.Addon, len(clusterAddons))
	for i, cA := range clusterAddons {
		if cA.Status == addonStatusReplaced {
			continue
		}
		if cA.Addon == nil {
			return nil, fmt.Errorf("addon was not specified for cluster addon %s", cA.ID)
		}
		addons[i] = *cA.Addon
	}
	return addons, nil
}

func filterAddons(addons []v1.Addon, name, version string) (r []v1.Addon) {
	for _, a := range addons {
		if (name == "" || a.Name == name) && (version == "" || a.ChartVersion == version) {
			r = append(r, a)
		}
	}
	return
}

func getConfigurationValues(addonID string, clusterAddons []clusteraddons.ClusterAddon) (string, error) {
	for _, a := range clusterAddons {
		if a.Addon.ID == addonID {
			return a.UserChartValues, nil
		}
	}
	return "", fmt.Errorf("addon %s was not found in specified cluster addons", addonID)
}

func removeInstalledAddons(availableAddons []v1.Addon, installedAddons []v1.Addon) (r []v1.Addon) {
OuterLoop:
	for _, av := range availableAddons {
		for _, in := range installedAddons {
			if av.ChartName == in.ChartName && av.ChartVersion == in.ChartVersion {
				continue OuterLoop
			}
		}
		r = append(r, av)
	}

	return
}
