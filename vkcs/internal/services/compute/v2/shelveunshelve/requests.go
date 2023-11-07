package shelveunshelve

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/shelveunshelve"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Shelve(client *gophercloud.ServiceClient, id string) shelveunshelve.ShelveResult {
	r := shelveunshelve.Shelve(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Unshelve(client *gophercloud.ServiceClient, id string, opts shelveunshelve.UnshelveOptsBuilder) shelveunshelve.UnshelveResult {
	r := shelveunshelve.Unshelve(client, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
