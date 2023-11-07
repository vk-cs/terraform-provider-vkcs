package floatingips

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Get(client *gophercloud.ServiceClient, id string) floatingips.GetResult {
	r := floatingips.Get(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func AssociateInstance(client *gophercloud.ServiceClient, serverID string, opts floatingips.AssociateOptsBuilder) floatingips.AssociateResult {
	r := floatingips.AssociateInstance(client, serverID, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func DisassociateInstance(client *gophercloud.ServiceClient, serverID string, opts floatingips.DisassociateOptsBuilder) floatingips.DisassociateResult {
	r := floatingips.DisassociateInstance(client, serverID, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
