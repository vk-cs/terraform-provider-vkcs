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

type UpdateOpts struct {
	ServerName string `json:"serverName"`
}

// Map builds the request body.
func (opts *UpdateOpts) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

type ProvisionOpts struct {
	v1.ProvisionFields
}

type UpdateNetworkConfigOpts struct {
	NetworkInterfaces []*v1.NetworkInterfaceConfig `json:"networkInterfaces,omitempty"`
	Bonds             []*v1.BondConfig             `json:"bonds,omitempty"`
}

func (opts *UpdateNetworkConfigOpts) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
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

// Update updates a baremetal server by ID.
func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(serverURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 204},
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

func UpdateNetworkConfig(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateNetworkConfigResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(updateNetworkConfigURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
