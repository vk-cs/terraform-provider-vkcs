package servers

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

type RenameOpts struct {
	ServerName string `json:"serverName"`
}

// Map builds the request body.
func (opts *RenameOpts) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

type ProvisionOpts struct {
	v1.ProvisionFields
}

// Map builds the request body.
func (opts *ProvisionOpts) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Get retrieves a baremetal server by ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(serverURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// Rename updates the name of a baremetal server by ID.
func Rename(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(serverURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// Delete removes a baremetal server by ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(serverURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// List returns a pager to iterate over all baremetal servers.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, serversURL(client), func(r pagination.PageResult) pagination.Page {
		return Page{
			TokenPageBase: paginationutil.TokenPageBase{PageResult: r},
		}
	})
}

// Provision starts provisioning of a baremetal server by ID.
func Provision(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r ProvisionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(provisionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
