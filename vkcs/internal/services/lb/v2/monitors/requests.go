package monitors

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/monitors"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(c *gophercloud.ServiceClient, opts monitors.CreateOptsBuilder) monitors.CreateResult {
	r := monitors.Create(c, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(c *gophercloud.ServiceClient, id string) monitors.GetResult {
	r := monitors.Get(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Update(c *gophercloud.ServiceClient, id string, opts monitors.UpdateOptsBuilder) monitors.UpdateResult {
	r := monitors.Update(c, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(c *gophercloud.ServiceClient, id string) monitors.DeleteResult {
	r := monitors.Delete(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
