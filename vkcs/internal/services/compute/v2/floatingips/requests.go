package floatingips

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func Get(client *gophercloud.ServiceClient, id string) floatingips.GetResult {
	r := floatingips.Get(client, id)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func AssociateInstance(client *gophercloud.ServiceClient, serverID string, opts floatingips.AssociateOptsBuilder) floatingips.AssociateResult {
	r := floatingips.AssociateInstance(client, serverID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}

func DisassociateInstance(client *gophercloud.ServiceClient, serverID string, opts floatingips.DisassociateOptsBuilder) floatingips.DisassociateResult {
	r := floatingips.DisassociateInstance(client, serverID, opts)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return r
}
