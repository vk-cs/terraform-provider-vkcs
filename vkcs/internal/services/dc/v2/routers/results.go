package routers

import (
	"github.com/gophercloud/gophercloud"
)

type RouterResp struct {
	Router Router `json:"dc_router"`
}

type Router struct {
	AvailabilityZone string `json:"availability_zone"`
	Flavor           string `json:"flavor"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	CreatedAt        string `json:"created_at"`
	ID               string `json:"id"`
	UpdatedAt        string `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*Router, error) {
	var res *RouterResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.Router, nil
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
