package images

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Get(client *gophercloud.ServiceClient, id string) images.GetResult {
	r := images.Get(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
