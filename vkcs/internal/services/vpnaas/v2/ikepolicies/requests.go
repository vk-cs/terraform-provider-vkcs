package ikepolicies

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ikepolicies"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Create(c *gophercloud.ServiceClient, opts ikepolicies.CreateOptsBuilder) ikepolicies.CreateResult {
	r := ikepolicies.Create(c, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Get(c *gophercloud.ServiceClient, id string) ikepolicies.GetResult {
	r := ikepolicies.Get(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(c *gophercloud.ServiceClient, id string) ikepolicies.DeleteResult {
	r := ikepolicies.Delete(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Update(c *gophercloud.ServiceClient, id string, opts ikepolicies.UpdateOptsBuilder) ikepolicies.UpdateResult {
	r := ikepolicies.Update(c, id, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
