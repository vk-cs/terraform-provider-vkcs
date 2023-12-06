package startstop

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/startstop"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func Start(client *gophercloud.ServiceClient, id string) startstop.StartResult {
	r := startstop.Start(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Stop(client *gophercloud.ServiceClient, id string) startstop.StopResult {
	r := startstop.Stop(client, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
