package secgroups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func AddServer(client *gophercloud.ServiceClient, serverID, groupName string) secgroups.AddServerResult {
	r := secgroups.AddServer(client, serverID, groupName)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func RemoveServer(client *gophercloud.ServiceClient, serverID, groupName string) secgroups.RemoveServerResult {
	r := secgroups.RemoveServer(client, serverID, groupName)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
