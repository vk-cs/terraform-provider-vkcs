package bgpstaticannounces

import (
	"github.com/gophercloud/gophercloud"
)

type BGPStaticAnnounceResp struct {
	BGPStaticAnnounce BGPStaticAnnounce `json:"dc_bgp_static_announce"`
}

type BGPStaticAnnounce struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DCBGPID     string `json:"dc_bgp_id"`
	Network     string `json:"network"`
	Gateway     string `json:"gateway"`
	Enabled     bool   `json:"enabled"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*BGPStaticAnnounce, error) {
	var res *BGPStaticAnnounceResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.BGPStaticAnnounce, nil
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
