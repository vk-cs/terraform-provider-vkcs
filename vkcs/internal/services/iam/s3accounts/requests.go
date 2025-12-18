package s3accounts

import (
	"github.com/gophercloud/gophercloud"
)

type CreateOptsBuilder interface {
	Map() (map[string]any, error)
}

// CreateOpts specifies attributes used to create a S3 account.
type CreateOpts struct {
	Name        string `json:"name" required:"true"`
	Description string `json:"description,omitempty"`
}

// Map builds a request body from a CreateOpts structure.
func (opts CreateOpts) Map() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Create implements a S3 account create request.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(s3AccountsURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get returns information about a S3 account, given its ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(s3AccountURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete implements a S3 account delete request.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(s3AccountURL(client, id), &gophercloud.RequestOpts{
		OkCodes:      []int{200, 204},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
