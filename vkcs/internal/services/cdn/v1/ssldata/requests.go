package ssldata

import "github.com/gophercloud/gophercloud"

type AddOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// AddOpts specifies attributes used to add a SSL certificate.
type AddOpts struct {
	Name           string `json:"name,omitempty"`
	SSLCertificate string `json:"sslCertificate,omitempty"`
	SSLPrivateKey  string `json:"sslPrivateKey,omitempty"`
}

// Map builds a request body from an AddOpts structure.
func (opts AddOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Add implements a SSL certificate add request.
func Add(client *gophercloud.ServiceClient, projectID string, opts AddOptsBuilder) (r AddResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(sslCertificatesURL(client, projectID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List returns a list of CDN SSL certificates in a project.
func List(client *gophercloud.ServiceClient, projectID string) (r ListResult) {
	resp, err := client.Get(sslCertificatesURL(client, projectID), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type UpdateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// UpdateOpts specifies attributes used to update a SSL certificate.
type UpdateOpts struct {
	Name string `json:"name,omitempty"`
}

// Map builds a request body from a UpdateOpts structure.
func (opts UpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Update implements a SSL certificate update request.
func Update(client *gophercloud.ServiceClient, projectID string, id int, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(sslCertificateURL(client, projectID, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete implements a resource delete request.
func Delete(client *gophercloud.ServiceClient, projectID string, id int) (r DeleteResult) {
	resp, err := client.Delete(sslCertificateURL(client, projectID, id), &gophercloud.RequestOpts{
		OkCodes:      []int{204},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
