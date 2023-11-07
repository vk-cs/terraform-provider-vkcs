package containers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/containers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Get(client *gophercloud.ServiceClient, id string) containers.GetResult {
	r := containers.Get(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Create(client *gophercloud.ServiceClient, opts containers.CreateOptsBuilder) containers.CreateResult {
	r := containers.Create(client, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(client *gophercloud.ServiceClient, id string) containers.DeleteResult {
	r := containers.Delete(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
