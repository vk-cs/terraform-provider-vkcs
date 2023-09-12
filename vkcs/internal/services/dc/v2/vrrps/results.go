package vrrps

import (
	"github.com/gophercloud/gophercloud"
)

type VRRPResp struct {
	VRRP VRRP `json:"dc_vrrp"`
}

type VRRP struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	GroupID        int    `json:"group_id"`
	NetworkID      string `json:"network_id"`
	SubnetID       string `json:"subnet_id"`
	AdvertInterval int    `json:"advert_interval"`
	Enabled        bool   `json:"enabled"`
	SDN            string `json:"sdn"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*VRRP, error) {
	var res *VRRPResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.VRRP, nil
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
