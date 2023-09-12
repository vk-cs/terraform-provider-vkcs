package vrrpaddresses

import (
	"github.com/gophercloud/gophercloud"
)

type VRRPAddressResp struct {
	VRRPAddress VRRPAddress `json:"dc_vrrp_address"`
}

type VRRPAddress struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DCVRRPID    string `json:"dc_vrrp_id"`
	IPAddress   string `json:"ip_address"`
	PortID      string `json:"port_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*VRRPAddress, error) {
	var res *VRRPAddressResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.VRRPAddress, nil
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
