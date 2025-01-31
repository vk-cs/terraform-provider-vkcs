package anycastips

import "github.com/gophercloud/gophercloud"

type CreateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// CreateOpts specifies attributes used to create a new anycast IP.
type CreateOpts struct {
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	NetworkID    string                 `json:"network_id" required:"true"`
	Associations []AnycastIPAssociation `json:"associations,omitempty"`
	HealthCheck  *AnycastIPHealthCheck  `json:"healthcheck,omitempty"`
}

// Map builds a request body from a CreateOpts structure.
func (opts CreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "anycastip")
	return b, err
}

// Create implements an anycast IP create request.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(anycastIPsURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get returns information about an anycast IP, given its ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(anycastIPURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type UpdateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// UpdateOpts specifies attributes used to update an anycast IP.
type UpdateOpts struct {
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Associations []AnycastIPAssociation `json:"associations,omitempty"`
	HealthCheck  *AnycastIPHealthCheck  `json:"healthcheck,omitempty"`
}

// Map builds a request body from a UpdateOpts structure.
func (opts UpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "anycastip")
	return b, err
}

// Update implements an anycast IP update request.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(anycastIPURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete implements an anycast IP delete request.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(anycastIPURL(client, id), &gophercloud.RequestOpts{
		OkCodes:      []int{204},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
