package vpnaas

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/endpointgroups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ikepolicies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/services"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/siteconnections"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
)

// EndpointGroupCreateOpts represents the attributes used when creating a new endpoint group.
type EndpointGroupCreateOpts struct {
	endpointgroups.CreateOpts
}

// IKEPolicyCreateOpts represents the attributes used when creating a new IKE policy.
type IKEPolicyCreateOpts struct {
	ikepolicies.CreateOpts
}

// IPSecPolicyCreateOpts represents the attributes used when creating a new IPSec policy.
type IPSecPolicyCreateOpts struct {
	ipsecpolicies.CreateOpts
}

// ServiceCreateOpts represents the attributes used when creating a new VPN service.
type ServiceCreateOpts struct {
	services.CreateOpts
}

// SiteConnectionCreateOpts represents the attributes used when creating a new IPSec site connection.
type SiteConnectionCreateOpts struct {
	siteconnections.CreateOpts
	TrafficSelectorEPMerge *bool `json:"traffic_selector_ep_merge,omitempty"`
}

// ToConnectionCreateMap casts a CreateOpts struct to a map.
func (opts SiteConnectionCreateOpts) ToConnectionCreateMap() (map[string]interface{}, error) {
	base, err := opts.CreateOpts.ToConnectionCreateMap()
	if err != nil {
		return nil, err
	}

	if opts.TrafficSelectorEPMerge == nil {
		return base, nil
	}

	siteConnection := base["ipsec_site_connection"].(map[string]interface{})
	siteConnection["traffic_selector_ep_merge"] = *opts.TrafficSelectorEPMerge

	return base, nil
}

type SiteConnectionUpdateOpts struct {
	siteconnections.UpdateOpts
	TrafficSelectorEPMerge *bool `json:"traffic_selector_ep_merge,omitempty"`
}

func (opts SiteConnectionUpdateOpts) ToConnectionUpdateMap() (map[string]interface{}, error) {
	base, err := opts.UpdateOpts.ToConnectionUpdateMap()
	if err != nil {
		return nil, err
	}

	if opts.TrafficSelectorEPMerge == nil {
		return base, nil
	}

	siteConnection := base["ipsec_site_connection"].(map[string]interface{})
	siteConnection["traffic_selector_ep_merge"] = *opts.TrafficSelectorEPMerge

	return base, nil
}

type groupExtended struct {
	endpointgroups.EndpointGroup
	networking.SDNExt
}

type ikePolicyExtended struct {
	ikepolicies.Policy
	networking.SDNExt
}

type ipsecPolicyExtended struct {
	ipsecpolicies.Policy
	networking.SDNExt
}

type serviceExtended struct {
	services.Service
	networking.SDNExt
}

type connectionExtended struct {
	siteconnections.Connection
	networking.SDNExt
	TrafficSelectorEPMerge bool `json:"traffic_selector_ep_merge,omitempty"`
}
