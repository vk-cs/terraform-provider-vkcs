package addons

import (
	"github.com/gophercloud/gophercloud"
)

type (
	CreateOpts struct {
		ClusterID      string
		AddonID        string `json:"addon_id"`
		AddonVersionID string `json:"addon_version_id"`
		Namespace      string `json:"namespace"`
		Values         string `json:"values"`
		AddonName      string `json:"addon_name"`
	}

	UpdateOpts struct {
		ClusterID      string
		ClusterAddonID string
		AddonVersionID string `json:"addon_version_id"`
		Values         string `json:"values"`
	}
)

func (opts CreateOpts) ToClusterAddonCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func (opts UpdateOpts) ToClusterAddonUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func ListAddons(client *gophercloud.ServiceClient) (res ListAddonsResult) {
	_, res.Err = client.Get(listURL(client), &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

func GetAddonVersionByName(client *gophercloud.ServiceClient, addonName, addonVersion string) (res GetAddonVersionResult) {
	_, res.Err = client.Get(getAddonByNameAndVersion(client, addonName, addonVersion), &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

func GetAddonVersionByID(client *gophercloud.ServiceClient, addonVersionID string) (res GetAddonVersionResult) {
	_, res.Err = client.Get(getAddonByGlobalID(client, addonVersionID), &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

func CreateClusterAddon(c *gophercloud.ServiceClient, opts *CreateOpts) CreateClusterAddonResult {
	var res CreateClusterAddonResult

	reqBody, err := opts.ToClusterAddonCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = c.Post(createClusterAddon(c, opts.ClusterID), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

func GetClusterAddon(c *gophercloud.ServiceClient, clusterAddonID string) GetClusterAddonResult {
	var res GetClusterAddonResult
	_, res.Err = c.Get(getClusterAddon(c, clusterAddonID), &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

func GetClusterAddonByClusterAndName(c *gophercloud.ServiceClient, clusterID, baseAddonName string) GetClusterAddonResult {
	var res GetClusterAddonResult
	_, res.Err = c.Get(getClusterAddonByClusterAndName(c, clusterID, baseAddonName), &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

func UpdateClusterAddon(c *gophercloud.ServiceClient, opts *UpdateOpts) error {
	reqBody, err := opts.ToClusterAddonUpdateMap()
	if err != nil {
		return err
	}

	_, err = c.Patch(updateClusterAddon(c, opts.ClusterID, opts.ClusterAddonID), reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return err
}

func DeleteClusterAddon(c *gophercloud.ServiceClient, clusterID, clusterAddonID string) error {
	_, err := c.Delete(deleteClusterAddon(c, clusterID, clusterAddonID), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return err
}
