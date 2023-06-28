package compute

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

const (
	antiAffinityPolicy = "anti-affinity"
	affinityPolicy     = "affinity"
)

// ServerGroupCreateOpts is a custom ServerGroup struct to include the
// ValueSpecs field.
type ComputeServerGroupCreateOpts struct {
	servergroups.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToServerGroupCreateMap casts a CreateOpts struct to a map.
// It overrides routers.ToServerGroupCreateMap to add the ValueSpecs field.
func (opts ComputeServerGroupCreateOpts) ToServerGroupCreateMap() (map[string]interface{}, error) {
	return util.BuildRequest(opts, "server_group")
}

func ExpandComputeServerGroupPolicies(client *gophercloud.ServiceClient, raw []interface{}) []string {
	policies := make([]string, len(raw))
	for i, v := range raw {
		client.Microversion = "2.15"

		policy := v.(string)
		policies[i] = policy

		// Set microversion for legacy policies to empty to not change behaviour for those policies
		if policy == antiAffinityPolicy || policy == affinityPolicy {
			client.Microversion = ""
		}
	}

	return policies
}
