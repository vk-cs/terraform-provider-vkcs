package s3accounts

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
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

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Offset and Limit are used
// for pagination.
type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

// ListOpts specifies options for listing S3 accounts.
type ListOpts struct {
	Offset         int    `q:"offset"`
	Limit          int    `q:"limit"`
	OrderBy        string `q:"order_by"`
	OrderDirection string `q:"order_direction"`
	Name           string `q:"name"`
}

// ToListQuery formats a ListOpts structure into a query string.
func (opts ListOpts) ToListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a paginated collection of S3 accounts.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) (r pagination.Pager) {
	url := s3AccountsURL(client)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return NewS3AccountPage(r)
	})
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
