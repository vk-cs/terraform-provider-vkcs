package tags

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/tags"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func List(client *gophercloud.ServiceClient, serverID string) tags.ListResult {
	r := tags.List(client, serverID)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func ReplaceAll(client *gophercloud.ServiceClient, serverID string, opts tags.ReplaceAllOptsBuilder) tags.ReplaceAllResult {
	r := tags.ReplaceAll(client, serverID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
