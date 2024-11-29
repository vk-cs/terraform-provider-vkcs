package attachinterfaces

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Get(client *gophercloud.ServiceClient, serverID, portID string) attachinterfaces.GetResult {
	r := attachinterfaces.Get(client, serverID, portID)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Create(client *gophercloud.ServiceClient, serverID string, opts attachinterfaces.CreateOptsBuilder) attachinterfaces.CreateResult {
	r := attachinterfaces.Create(client, serverID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, serverID, portID string) attachinterfaces.DeleteResult {
	r := attachinterfaces.Delete(client, serverID, portID)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
