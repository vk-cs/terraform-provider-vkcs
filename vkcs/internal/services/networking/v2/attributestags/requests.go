package attributestags

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func ReplaceAll(client *gophercloud.ServiceClient, resourceType string, resourceID string, opts attributestags.ReplaceAllOptsBuilder) attributestags.ReplaceAllResult {
	r := attributestags.ReplaceAll(client, resourceType, resourceID, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
