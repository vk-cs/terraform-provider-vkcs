package addons

import (
	"github.com/gophercloud/gophercloud"
)

type (
	AddonList struct {
		Addons []Addon `json:"addons"`
	}

	Addon struct {
		ID       string              `json:"id"`
		Name     string              `json:"name"`
		Versions []AddonVersionShort `json:"versions"`
	}

	AddonVersionShort struct {
		ID           string `json:"id"`
		Version      string `json:"version"`
		IsDeprecated bool   `json:"is_deprecated"`
	}

	AddonVersion struct {
		ID                    string   `json:"id"`
		AddonID               string   `json:"addon_id"`
		Name                  string   `json:"name"`
		Version               string   `json:"version"`
		ValuesTemplate        *string  `json:"values_template"`
		SupportedKubeVersions []string `json:"supported_kube_versions"`
	}

	ClusterAddonID struct {
		ID string `json:"id"`
	}

	ClusterAddon struct {
		ID             string `json:"id"`
		ClusterID      string `json:"cluster_id"`
		AddonID        string `json:"addon_id"`
		AddonVersionID string `json:"addon_version_id"`
		Namespace      string `json:"namespace"`
		Values         string `json:"values"`
		Status         string `json:"status"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at"`
		AddonName      string `json:"addon_name"`
		BaseAddonName  string `json:"base_addon_name"`
	}

	ListAddonsResult struct {
		gophercloud.Result
	}

	GetAddonVersionResult struct {
		gophercloud.Result
	}

	CreateClusterAddonResult struct {
		gophercloud.Result
	}

	GetClusterAddonResult struct {
		gophercloud.Result
	}
)

func (r ListAddonsResult) Extract() (AddonList, error) {
	if r.Err != nil {
		return AddonList{}, r.Err
	}

	var addonList AddonList
	err := r.ExtractInto(&addonList)
	return addonList, err
}

func (r GetAddonVersionResult) Extract() (AddonVersion, error) {
	if r.Err != nil {
		return AddonVersion{}, r.Err
	}

	var addonVersion AddonVersion
	err := r.ExtractInto(&addonVersion)
	return addonVersion, err
}

func (r CreateClusterAddonResult) Extract() (ClusterAddonID, error) {
	if r.Err != nil {
		return ClusterAddonID{}, r.Err
	}

	var clusterAddonID ClusterAddonID
	err := r.ExtractInto(&clusterAddonID)
	return clusterAddonID, err
}

func (r GetClusterAddonResult) Extract() (ClusterAddon, error) {
	if r.Err != nil {
		return ClusterAddon{}, r.Err
	}

	var clusterAddon ClusterAddon
	err := r.ExtractInto(&clusterAddon)
	return clusterAddon, err
}
