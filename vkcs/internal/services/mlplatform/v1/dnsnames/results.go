package dnsnames

import (
	"github.com/gophercloud/gophercloud"
)

type Response struct {
	DNS string `json:"dns" required:"true"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*Response, error) {
	var res *Response
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return res, nil
}

type GetResult struct {
	commonResult
}
