package zones

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type CreateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// CreateOpts specifies the attributes used to create a zone.
type CreateOpts struct {
	Zone          string `json:"zone" required:"true"`
	SOAPrimaryDNS string `json:"soa_primary_dns,omitempty"`
	SOAAdminEmail string `json:"soa_admin_email,omitempty"`
	SOARefresh    int    `json:"soa_refresh,omitempty"`
	SOARetry      int    `json:"soa_retry,omitempty"`
	SOAExpire     int    `json:"soa_expire,omitempty"`
	SOATTL        int    `json:"soa_ttl,omitempty"`
}

// Map formats a CreateOpts structure into a request body.
func (opts CreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Create implements a zone create request.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(zonesURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get returns information about a zone, given its ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(zoneURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// list request.
type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned.
type ListOpts struct {
	// Integer value for the limit of values to return.
	Limit int `q:"limit"`

	// UUID of the zone at which you want to set a marker.
	Marker string `q:"marker"`

	ID            string `q:"uuid"`
	Zone          string `q:"zone"`
	SOAPrimaryDNS string `q:"soa_primary_dns"`
	SOAAdminEmail string `q:"soa_admin_email"`
	SOASerial     int    `q:"soa_serial"`
	SOARefresh    int    `q:"soa_refresh"`
	SOARetry      int    `q:"soa_retry"`
	SOAExpire     int    `q:"soa_expire"`
	SOATTL        int    `q:"soa_ttl"`
	Status        string `q:"status"`
}

// ToListQuery formats a listOpts structure into a query string.
func (opts ListOpts) ToListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List implements a zone list request.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) (r ListResult) {
	url := zonesURL(client)
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Get(url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

type UpdateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// UpdateOpts specifies the attributes to update a zone.
type UpdateOpts struct {
	SOAPrimaryDNS string `json:"soa_primary_dns,omitempty"`
	SOAAdminEmail string `json:"soa_admin_email,omitempty"`
	SOARefresh    int    `json:"soa_refresh,omitempty"`
	SOARetry      int    `json:"soa_retry,omitempty"`
	SOAExpire     int    `json:"soa_expire,omitempty"`
	SOATTL        int    `json:"soa_ttl,omitempty"`
}

// Map formats a zoneUpdateOpts structure into a request body.
func (opts UpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Update implements a zone update request.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(zoneURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete implements a zone delete request.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(zoneURL(client, id), &gophercloud.RequestOpts{
		OkCodes:      []int{204},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
