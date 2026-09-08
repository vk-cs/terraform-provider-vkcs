package rents

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	paginationutil "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/pagination"
)

type RentRequest struct {
	RentRequestId string   `json:"rentRequestId"`
	ServerIds     []string `json:"serverIds"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*RentRequest, error) {
	var rent RentRequest
	err := r.ExtractInto(&rent)
	return &rent, err
}

// CreateResult represents result of baremetal rent request create.
type CreateResult struct {
	commonResult
}

// GetResult represents result of baremetal rent request get.
type GetResult struct {
	commonResult
}

// Page represents a page of baremetal rent request.
type Page struct {
	paginationutil.TokenPageBase
}

func ExtractRentRequests(p pagination.Page) ([]RentRequest, error) {
	var s struct {
		Items []RentRequest `json:"items"`
	}
	err := p.(Page).ExtractInto(&s)
	return s.Items, err
}
