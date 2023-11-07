package volumeattach

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(client *gophercloud.ServiceClient, serverID string, opts volumeattach.CreateOptsBuilder) volumeattach.CreateResult {
	r := volumeattach.Create(client, serverID, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(client *gophercloud.ServiceClient, serverID, attachmentID string) volumeattach.GetResult {
	r := volumeattach.Get(client, serverID, attachmentID)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, serverID, attachmentID string) volumeattach.DeleteResult {
	r := volumeattach.Delete(client, serverID, attachmentID)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
