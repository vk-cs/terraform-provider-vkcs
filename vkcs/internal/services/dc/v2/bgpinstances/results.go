package bgpinstances

import (
	"github.com/gophercloud/gophercloud"
)

type BGPInstanceResp struct {
	BGPInstance BGPInstance `json:"dc_bgp"`
}

type BGPInstance struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	Description              string `json:"description"`
	DCRouterID               string `json:"dc_router_id"`
	BGPRouterID              string `json:"bgp_router_id"`
	ASN                      int    `json:"asn"`
	ECMPEnabled              bool   `json:"ecmp_enabled"`
	Enabled                  bool   `json:"enabled"`
	GracefulRestart          bool   `json:"graceful_restart"`
	LongLivedGracefulRestart bool   `json:"long_lived_graceful_restart"`
	CreatedAt                string `json:"created_at"`
	UpdatedAt                string `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*BGPInstance, error) {
	var res *BGPInstanceResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.BGPInstance, nil
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
