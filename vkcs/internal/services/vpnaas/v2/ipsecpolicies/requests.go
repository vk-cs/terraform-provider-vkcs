package ipsecpolicies

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(c *gophercloud.ServiceClient, opts ipsecpolicies.CreateOptsBuilder) ipsecpolicies.CreateResult {
	r := ipsecpolicies.Create(c, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(c *gophercloud.ServiceClient, id string) ipsecpolicies.DeleteResult {
	r := ipsecpolicies.Delete(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(c *gophercloud.ServiceClient, id string) ipsecpolicies.GetResult {
	r := ipsecpolicies.Get(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Update(c *gophercloud.ServiceClient, id string, opts ipsecpolicies.UpdateOptsBuilder) ipsecpolicies.UpdateResult {
	r := ipsecpolicies.Update(c, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
