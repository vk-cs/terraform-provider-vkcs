package addons

import (
	"encoding/base64"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"

	v1 "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfraaddons/v1"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type InstallAddonToClusterOpts struct {
	Values   string     `json:"values"`
	Payload  v1.Payload `json:"payload"`
	Migrated bool       `json:"migrated"`
}

func (opts InstallAddonToClusterOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.Values != "" {
		valuesBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(opts.Values)))
		base64.StdEncoding.Encode(valuesBase64, []byte(opts.Values))
		b["values"] = string(valuesBase64)
	}

	return b, nil
}

func GetAvailableAddon(client *gophercloud.ServiceClient, clusterID, addonID string) (r GetAvailableAddonResult) {
	resp, err := client.Get(clusterAvailableAddonURL(client, clusterID, addonID), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func ListClusterAddons(client *gophercloud.ServiceClient, clusterID string) pagination.Pager {
	return pagination.NewPager(client, clusterAddonsURL(client, clusterID), func(r pagination.PageResult) pagination.Page {
		return ClusterAddonPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func ListClusterAvailableAddons(client *gophercloud.ServiceClient, clusterID string) pagination.Pager {
	return pagination.NewPager(client, clusterAvailableAddonsURL(client, clusterID), func(r pagination.PageResult) pagination.Page {
		return AddonPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func InstallAddonToCluster(client *gophercloud.ServiceClient, addonID, clusterID string, opts OptsBuilder) (r InstallAddonToClusterResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(installAddonToClusterURL(client, addonID, clusterID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
