package apioptions

import (
	"github.com/gophercloud/gophercloud"
)

type APIOptionsResp struct {
	APIOptions APIOptions `json:"dc_api_options"`
}

type APIOptions struct {
	AvailabilityZones []string `json:"availability_zones"`
	Flavors           []string `json:"flavors"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*APIOptions, error) {
	var res *APIOptionsResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.APIOptions, nil
}

type GetResult struct {
	commonResult
}
