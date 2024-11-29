package secgroups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func AddServer(client *gophercloud.ServiceClient, serverID, groupName string) secgroups.AddServerResult {
	r := secgroups.AddServer(client, serverID, groupName)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func RemoveServer(client *gophercloud.ServiceClient, serverID, groupName string) secgroups.RemoveServerResult {
	r := secgroups.RemoveServer(client, serverID, groupName)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
