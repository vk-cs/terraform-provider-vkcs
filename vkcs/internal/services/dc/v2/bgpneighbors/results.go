package bgpneighbors

import (
	"github.com/gophercloud/gophercloud"
)

type BGPNeighborResp struct {
	BGPNeighbor BGPNeighbor `json:"dc_bgp_neighbor"`
}

type BGPNeighbor struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	DCBGPID              string `json:"dc_bgp_id"`
	RemoteASN            int    `json:"remote_asn"`
	RemoteIP             string `json:"remote_ip"`
	ForceIBGPNextHopSelf bool   `json:"force_ibgp_next_hop_self"`
	AddPaths             string `json:"add_paths"`
	BFDEnabled           bool   `json:"bfd_enabled"`
	FilterIn             string `json:"filter_in"`
	FilterOut            string `json:"filter_out"`
	Enabled              bool   `json:"enabled"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*BGPNeighbor, error) {
	var res *BGPNeighborResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.BGPNeighbor, nil
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}
