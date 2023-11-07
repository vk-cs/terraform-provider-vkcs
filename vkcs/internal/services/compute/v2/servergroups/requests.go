package servergroups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(client *gophercloud.ServiceClient, opts servergroups.CreateOptsBuilder) servergroups.CreateResult {
	r := servergroups.Create(client, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(client *gophercloud.ServiceClient, id string) servergroups.GetResult {
	r := servergroups.Get(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, id string) servergroups.DeleteResult {
	r := servergroups.Delete(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
