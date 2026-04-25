package rents

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	v1 "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/baremetal/v1"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	paginationutil "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/pagination"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type CreateOpts struct {
	v1.ProvisionFields
	FlavorId         string   `json:"flavorId"`
	ServerCount      int64    `json:"serverCount"`
	AvailabilityZone *string  `json:"availabilityZone"`
	Tags             []string `json:"tags"`
}

// Map builds the request body.
func (opts *CreateOpts) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create sends a request to create a baremetal rent request.
func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(rentRequestsURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// Get retrieves a baremetal rent request by ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(rentRequestURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// List returns a pager to iterate over all baremetal rent requests.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	url := rentRequestsURL(client)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return Page{
			TokenPageBase: paginationutil.TokenPageBase{PageResult: r},
		}
	})
}
