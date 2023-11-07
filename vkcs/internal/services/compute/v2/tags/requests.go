package tags

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/tags"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func List(client *gophercloud.ServiceClient, serverID string) tags.ListResult {
	r := tags.List(client, serverID)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func ReplaceAll(client *gophercloud.ServiceClient, serverID string, opts tags.ReplaceAllOptsBuilder) tags.ReplaceAllResult {
	r := tags.ReplaceAll(client, serverID, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
