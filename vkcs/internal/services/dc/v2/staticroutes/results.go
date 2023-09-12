package staticroutes

import (
	"github.com/gophercloud/gophercloud"
)

type StaticRouteResp struct {
	StaticRoute StaticRoute `json:"dc_static_route"`
}

type StaticRoute struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DCRouterID  string `json:"dc_router_id"`
	Network     string `json:"network"`
	Gateway     string `json:"gateway"`
	Metric      int    `json:"metric"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*StaticRoute, error) {
	var res *StaticRouteResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.StaticRoute, nil
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
