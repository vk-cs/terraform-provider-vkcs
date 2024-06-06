package conntrackhelpers

import (
	"github.com/gophercloud/gophercloud"
)

type ConntrackHelperResp struct {
	ConntrackHelper ConntrackHelper `json:"dc_conntrack_helper"`
}

type ConntrackHelper struct {
	ID          string `json:"id"`
	DCRouterID  string `json:"dc_router_id"`
	Protocol    string `json:"protocol"`
	Port        int    `json:"port"`
	Helper      string `json:"helper"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*ConntrackHelper, error) {
	var res *ConntrackHelperResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.ConntrackHelper, nil
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
