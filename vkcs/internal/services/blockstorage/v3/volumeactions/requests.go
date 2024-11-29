package volumeactions

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func ExtendSize(client *gophercloud.ServiceClient, id string, opts volumeactions.ExtendSizeOptsBuilder) volumeactions.ExtendSizeResult {
	r := volumeactions.ExtendSize(client, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func ChangeType(client *gophercloud.ServiceClient, id string, opts volumeactions.ChangeTypeOptsBuilder) volumeactions.ChangeTypeResult {
	r := volumeactions.ChangeType(client, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
