package bgpneighbors

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type BGPNeighborCreate struct {
	BGPNeighbor *CreateOpts `json:"dc_bgp_neighbor"`
}

type CreateOpts struct {
	Name                 string `json:"name,omitempty"`
	Description          string `json:"description,omitempty"`
	DCBGPID              string `json:"dc_bgp_id"`
	RemoteASN            int    `json:"remote_asn"`
	RemoteIP             string `json:"remote_ip"`
	ForceIBGPNextHopSelf *bool  `json:"force_ibgp_next_hop_self,omitempty"`
	AddPaths             string `json:"add_paths,omitempty"`
	BFDEnabled           *bool  `json:"bfd_enabled,omitempty"`
	FilterIn             string `json:"filter_in,omitempty"`
	FilterOut            string `json:"filter_out,omitempty"`
	Enabled              *bool  `json:"enabled,omitempty"`
}

func (opts *BGPNeighborCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(bgpNeighborsURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(bgpNeighborURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

type BGPNeighborUpdate struct {
	BGPNeighbor *UpdateOpts `json:"dc_bgp_neighbor"`
}

type UpdateOpts struct {
	Name                 string `json:"name,omitempty"`
	Description          string `json:"description,omitempty"`
	ForceIBGPNextHopSelf *bool  `json:"force_ibpg_next_hop_self,omitempty"`
	AddPaths             string `json:"add_paths,omitempty"`
	BFDEnabled           *bool  `json:"bfd_enabled,omitempty"`
	FilterIn             string `json:"filter_in,omitempty"`
	FilterOut            string `json:"filter_out,omitempty"`
	Enabled              *bool  `json:"enabled,omitempty"`
}

func (opts *BGPNeighborUpdate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(bgpNeighborURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(bgpNeighborURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
