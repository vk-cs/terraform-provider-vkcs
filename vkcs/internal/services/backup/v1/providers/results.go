package providers

import (
	"github.com/gophercloud/gophercloud"
)

type ProviderResp struct {
	Provider Provider `json:"provider"`
}

type ProvidersResp struct {
	Providers []*Provider `json:"providers"`
}

type Provider struct {
	ID   string `json:"id" required:"true"`
	Name string `json:"name" required:"true"`
}

type ListResult struct {
	gophercloud.Result
}

func (r ListResult) Extract() ([]*Provider, error) {
	var s *ProvidersResp
	if err := r.ExtractInto(&s); err != nil {
		return nil, err
	}
	return s.Providers, nil
}
