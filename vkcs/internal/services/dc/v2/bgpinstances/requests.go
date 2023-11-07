package bgpinstances

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type BGPInstanceCreate struct {
	BGPInstance *CreateOpts `json:"dc_bgp"`
}

type CreateOpts struct {
	Name                     string `json:"name,omitempty"`
	Description              string `json:"description,omitempty"`
	DCRouterID               string `json:"dc_router_id"`
	BGPRouterID              string `json:"bgp_router_id"`
	ASN                      int    `json:"asn"`
	ECMPEnabled              *bool  `json:"ecmp_enabled,omitempty"`
	Enabled                  *bool  `json:"enabled,omitempty"`
	GracefulRestart          *bool  `json:"graceful_restart,omitempty"`
	LongLivedGracefulRestart *bool  `json:"long_lived_graceful_restart,omitempty"`
}

func (opts *BGPInstanceCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(bgpsURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(bgpURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}

type BGPInstanceUpdate struct {
	BGPInstance *UpdateOpts `json:"dc_bgp"`
}

type UpdateOpts struct {
	Name                     string `json:"name,omitempty"`
	Description              string `json:"description,omitempty"`
	ECMPEnabled              *bool  `json:"ecmp_enabled,omitempty"`
	Enabled                  *bool  `json:"enabled,omitempty"`
	GracefulRestart          *bool  `json:"graceful_restart,omitempty"`
	LongLivedGracefulRestart *bool  `json:"long_lived_graceful_restart,omitempty"`
}

func (opts *BGPInstanceUpdate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(bgpURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(bgpURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}
