package flavors

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Get(client *gophercloud.ServiceClient, id string) flavors.GetResult {
	r := flavors.Get(client, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func ListExtraSpecs(client *gophercloud.ServiceClient, flavorID string) flavors.ListExtraSpecsResult {
	r := flavors.ListExtraSpecs(client, flavorID)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
