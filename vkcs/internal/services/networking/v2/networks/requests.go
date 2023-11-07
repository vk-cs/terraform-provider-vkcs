package networks

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

// NetworkCreateOpts represents the attributes used when creating a new network.
type NetworkCreateOpts struct {
	networks.CreateOpts
	ValueSpecs       map[string]string `json:"value_specs,omitempty"`
	PrivateDNSDomain string            `json:"private_dns_domain,omitempty"`
	ServicesAccess   bool              `json:"enable_shadow_port,omitempty"`
}

// ToNetworkCreateMap casts a CreateOpts struct to a map.
// It overrides networks.ToNetworkCreateMap to add the ValueSpecs field.
func (opts NetworkCreateOpts) ToNetworkCreateMap() (map[string]interface{}, error) {
	return util.BuildRequest(opts, "network")
}

// NetworkUpdateOpts represents the attributes used when updating a network.
type NetworkUpdateOpts struct {
	networks.UpdateOpts
	ServicesAccess *bool `json:"enable_shadow_port,omitempty"`
}

// ToNetworkUpdateMap casts a UpdateOpts struct to a map.
// It overrides networks.ToNetworkUpdateMap to add the ServicesAccess field.
func (opts NetworkUpdateOpts) ToNetworkUpdateMap() (map[string]interface{}, error) {
	return util.BuildRequest(opts, "network")
}

func Get(c *gophercloud.ServiceClient, id string) networks.GetResult {
	r := networks.Get(c, id)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Create(c *gophercloud.ServiceClient, opts networks.CreateOptsBuilder) networks.CreateResult {
	r := networks.Create(c, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Update(c *gophercloud.ServiceClient, networkID string, opts networks.UpdateOptsBuilder) networks.UpdateResult {
	r := networks.Update(c, networkID, opts)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}

func Delete(c *gophercloud.ServiceClient, networkID string) networks.DeleteResult {
	r := networks.Delete(c, networkID)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return r
}
