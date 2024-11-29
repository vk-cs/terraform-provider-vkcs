package services

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/services"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Create(c *gophercloud.ServiceClient, opts services.CreateOptsBuilder) services.CreateResult {
	r := services.Create(c, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Delete(c *gophercloud.ServiceClient, id string) services.DeleteResult {
	r := services.Delete(c, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Update(c *gophercloud.ServiceClient, id string, opts services.UpdateOptsBuilder) services.UpdateResult {
	r := services.Update(c, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Get(c *gophercloud.ServiceClient, id string) services.GetResult {
	r := services.Get(c, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
