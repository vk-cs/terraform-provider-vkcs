package listeners

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Create(c *gophercloud.ServiceClient, opts listeners.CreateOptsBuilder) listeners.CreateResult {
	r := listeners.Create(c, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Get(c *gophercloud.ServiceClient, id string) listeners.GetResult {
	r := listeners.Get(c, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Update(c *gophercloud.ServiceClient, id string, opts listeners.UpdateOpts) listeners.UpdateResult {
	r := listeners.Update(c, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Delete(c *gophercloud.ServiceClient, id string) listeners.DeleteResult {
	r := listeners.Delete(c, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
