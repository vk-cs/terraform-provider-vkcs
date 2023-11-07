package rules

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(c *gophercloud.ServiceClient, opts rules.CreateOptsBuilder) rules.CreateResult {
	r := rules.Create(c, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(c *gophercloud.ServiceClient, id string) rules.GetResult {
	r := rules.Get(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(c *gophercloud.ServiceClient, id string) rules.DeleteResult {
	r := rules.Delete(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
