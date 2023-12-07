package imageinfos

import (
	"github.com/gophercloud/gophercloud"
)

type Response struct {
	VolumeSize int64 `json:"volume_size" required:"true"`
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
