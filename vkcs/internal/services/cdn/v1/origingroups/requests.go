package origingroups

import "github.com/gophercloud/gophercloud"

type CreateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// CreateOpts specifies attributes used to create a CDN origin group.
type CreateOpts struct {
	Name    string   `json:"name"`
	Origins []Origin `json:"origins"`
	UseNext bool     `json:"useNext"`
}

// Map builds a request body from a CreateOpts structure.
func (opts CreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Create implements an origin group create request.
func Create(client *gophercloud.ServiceClient, projectID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(originGroupsURL(client, projectID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get returns information about an origin group, given its ID.
func Get(client *gophercloud.ServiceClient, projectID string, id int) (r GetResult) {
	resp, err := client.Get(originGroupURL(client, projectID, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List returns a list of origin groups in a project.
func List(client *gophercloud.ServiceClient, projectID string) (r ListResult) {
	resp, err := client.Get(originGroupsURL(client, projectID), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type UpdateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// UpdateOpts specifies attributes used to update a CDN origin group.
type UpdateOpts struct {
	Name    string   `json:"name"`
	Origins []Origin `json:"origins,omitempty"`
	UseNext bool     `json:"useNext"`
}

// Map builds a request body from a UpdateOpts structure.
func (opts UpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Update implements an origin group update request.
func Update(client *gophercloud.ServiceClient, projectID string, id int, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(originGroupURL(client, projectID, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete implements an origin group delete request.
func Delete(client *gophercloud.ServiceClient, projectID string, id int) (r DeleteResult) {
	resp, err := client.Delete(originGroupURL(client, projectID, id), &gophercloud.RequestOpts{
		OkCodes:      []int{204},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
