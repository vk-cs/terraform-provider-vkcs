package regions

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/regions"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Get(client *gophercloud.ServiceClient, id string) regions.GetResult {
	r := regions.Get(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
