package addons

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"

	v1 "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfraaddons/v1"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfraaddons/v1/clusteraddons"
)

type GetAvailableAddonResult struct {
	gophercloud.Result
}

func (r GetAvailableAddonResult) Extract() (*v1.Addon, error) {
	var s struct {
		ClusterID string    `json:"cluster_id"`
		Addon     *v1.Addon `json:"addon"`
	}
	err := r.ExtractInto(&s)
	return s.Addon, err
}

type InstallAddonToClusterResult struct {
	gophercloud.Result
}

func (r InstallAddonToClusterResult) Extract() (*clusteraddons.ClusterAddon, error) {
	var s struct {
		ClusterAddon *clusteraddons.ClusterAddon `json:"addon"`
	}
	err := r.ExtractInto(&s)
	return s.ClusterAddon, err
}

type AddonPage struct {
	pagination.LinkedPageBase
}

func (r AddonPage) IsEmpty() (bool, error) {
	s, err := ExtractAddons(r)
	return len(s) == 0, err
}

func ExtractAddons(r pagination.Page) ([]v1.Addon, error) {
	var s struct {
		ClusterID string     `json:"cluster_id"`
		Addons    []v1.Addon `json:"addons"`
	}
	err := (r.(AddonPage)).ExtractInto(&s)
	return s.Addons, err
}

type ClusterAddonPage struct {
	pagination.LinkedPageBase
}

func (r ClusterAddonPage) IsEmpty() (bool, error) {
	s, err := ExtractClusterAddons(r)
	return len(s) == 0, err
}

func ExtractClusterAddons(r pagination.Page) ([]clusteraddons.ClusterAddon, error) {
	var s struct {
		ClusterID string                       `json:"cluster_id"`
		Addons    []clusteraddons.ClusterAddon `json:"addons"`
	}
	err := (r.(ClusterAddonPage)).ExtractInto(&s)
	return s.Addons, err
}
