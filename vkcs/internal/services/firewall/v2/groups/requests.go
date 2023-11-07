package groups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(c *gophercloud.ServiceClient, opts groups.CreateOptsBuilder) groups.CreateResult {
	r := groups.Create(c, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Update(c *gophercloud.ServiceClient, id string, opts groups.UpdateOptsBuilder) groups.UpdateResult {
	r := groups.Update(c, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(c *gophercloud.ServiceClient, id string) groups.GetResult {
	r := groups.Get(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(c *gophercloud.ServiceClient, id string) groups.DeleteResult {
	r := groups.Delete(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
