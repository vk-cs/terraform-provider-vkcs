package vrrpinterfaces

import (
	"github.com/gophercloud/gophercloud"
)

type VRRPInterfaceResp struct {
	VRRPInterface VRRPInterface `json:"dc_vrrp_interface"`
}

type VRRPInterface struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	DCVRRPID      string `json:"dc_vrrp_id"`
	DCInterfaceID string `json:"dc_interface_id"`
	Priority      int    `json:"priority"`
	Preempt       bool   `json:"preempt"`
	Master        bool   `json:"master"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*VRRPInterface, error) {
	var res *VRRPInterfaceResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.VRRPInterface, nil
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
