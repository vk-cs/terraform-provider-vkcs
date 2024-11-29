package siteconnections

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/siteconnections"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Create(c *gophercloud.ServiceClient, opts siteconnections.CreateOptsBuilder) siteconnections.CreateResult {
	r := siteconnections.Create(c, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Delete(c *gophercloud.ServiceClient, id string) siteconnections.DeleteResult {
	r := siteconnections.Delete(c, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Get(c *gophercloud.ServiceClient, id string) siteconnections.GetResult {
	r := siteconnections.Get(c, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func Update(c *gophercloud.ServiceClient, id string, opts siteconnections.UpdateOptsBuilder) siteconnections.UpdateResult {
	r := siteconnections.Update(c, id, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
