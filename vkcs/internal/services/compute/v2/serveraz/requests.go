package serveraz

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type ChangeOpts struct {
	AvailabilityZone string `json:"availability_zone"`
}

func (opts ChangeOpts) ToChangeOptsMap() (map[string]interface{}, error) {
	return util.BuildRequest(opts, "os-changeAvailabilityZone")
}

// ChangeAvailabilityZone changes the availability zone of a server using the
// os-changeAvailabilityZone server action (microversion 2.42). The operation
// is asynchronous: an ACTIVE or PAUSED server is live-migrated, a STOPPED,
// SHELVED, SHELVED_OFFLOADED or STOPPED_OFFLOADED server is cold-migrated.
func ChangeAvailabilityZone(client *gophercloud.ServiceClient, id string, opts ChangeOpts) (r gophercloud.ErrResult) {
	b, err := opts.ToChangeOptsMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(client.ServiceURL("servers", id, "action"), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
