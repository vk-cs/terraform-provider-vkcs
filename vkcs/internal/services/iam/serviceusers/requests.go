package serviceusers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type CreateOptsBuilder interface {
	Map() (map[string]any, error)
}

// CreateOpts specifies attributes used to create a service user.
type CreateOpts struct {
	Name        string   `json:"name" required:"true"`
	RoleNames   []string `json:"role_names" required:"true"`
	Description string   `json:"description,omitempty"`
}

// Map builds a request body from a CreateOpts structure.
func (opts CreateOpts) Map() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Create implements a service user create request.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(serviceUsersURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get returns information about a service user, given its ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(serviceUserURL(client, id), &r.Body, &gophercloud.RequestOpts{
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

// ListOpts specifies options for listing service users.
type ListOpts struct {
	Limit          int    `q:"limit"`
	Offset         int    `q:"offset"`
	OrderBy        string `q:"order_by"`
	OrderDirection string `q:"order_direction"`
	Name           string `q:"name"`
}

// ToListQuery formats a listOpts structure into a query string.
func (opts ListOpts) ToListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a paginated collection of service users.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) (r pagination.Pager) {
	url := serviceUsersURL(client)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ServiceUserPage{PageResult: r}
	})
}

// Delete implements a service user delete request.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(serviceUserURL(client, id), &gophercloud.RequestOpts{
		OkCodes:      []int{200, 204},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
