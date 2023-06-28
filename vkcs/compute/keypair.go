package compute

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

// ComputeKeyPairV2CreateOpts is a custom KeyPair struct to include the ValueSpecs field.
type ComputeKeyPairV2CreateOpts struct {
	keypairs.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToKeyPairCreateMap casts a CreateOpts struct to a map.
// It overrides keypairs.ToKeyPairCreateMap to add the ValueSpecs field.
func (opts ComputeKeyPairV2CreateOpts) ToKeyPairCreateMap() (map[string]interface{}, error) {
	return util.BuildRequest(opts, "keypair")
}
