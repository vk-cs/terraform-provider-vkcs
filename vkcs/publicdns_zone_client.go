package vkcs

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

type publicDNSClient interface {
	Get(url string, JSONResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	Post(url string, JSONBody interface{}, JSONResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	Patch(url string, JSONBody interface{}, JSONResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	Delete(url string, opts *gophercloud.RequestOpts) (*http.Response, error)
	Head(url string, opts *gophercloud.RequestOpts) (*http.Response, error)
	Put(url string, JSONBody interface{}, JSONResponse interface{}, opts *gophercloud.RequestOpts) (*http.Response, error)
	ServiceURL(parts ...string) string
}

type commonZoneResult struct {
	gophercloud.Result
}

// Extract interprets a zoneGetResult, zoneCreateResult or zoneUpdateResult as a zone.
// An error is returned if the original call or the extraction failed.
func (r commonZoneResult) Extract() (*zone, error) {
	var s *zone
	err := r.ExtractInto(&s)
	return s, err
}

// zoneCreateResult is the result of a zoneCreate request. Call its Extract method
// to interpret the result as a zone.
type zoneCreateResult struct {
	commonZoneResult
}

// zoneGetResult is the result of a zoneGet request. Call its Extract method
// to interpret the result as a zone.
type zoneGetResult struct {
	commonZoneResult
}

// zoneListResult is the result of a zoneList request. Call its Extract method
// to interpret the result as a slice of zones.
type zoneListResult struct {
	commonZoneResult
}

// Extract extracts a slice of zones from a zoneListResult.
func (r zoneListResult) Extract() ([]zone, error) {
	var zones []zone
	err := r.ExtractInto(&zones)
	return zones, err
}

// zoneUpdateResult is the result of a zoneUpdate request. Call its Extract method
// to interpret the result as a zone.
type zoneUpdateResult struct {
	commonZoneResult
}

// zoneDeleteResult is the result of a zoneDelete request. Call its ExtractErr method
// to determine if the request succeeded or failed.
type zoneDeleteResult struct {
	gophercloud.ErrResult
}

// zone represents a public DNS zone.
type zone struct {
	ID         string `json:"uuid"`
	Zone       string `json:"zone"`
	PrimaryDNS string `json:"soa_primary_dns"`
	AdminEmail string `json:"soa_admin_email"`
	Serial     int    `json:"soa_serial"`
	Refresh    int    `json:"soa_refresh"`
	Retry      int    `json:"soa_retry"`
	Expire     int    `json:"soa_expire"`
	TTL        int    `json:"soa_ttl"`
	Status     string `json:"status"`
}

var zonesAPIPath = "dns"

// zoneCreateOpts specifies the attributes used to create a zone.
type zoneCreateOpts struct {
	Zone       string `json:"zone" required:"true"`
	PrimaryDNS string `json:"soa_primary_dns,omitempty"`
	AdminEmail string `json:"soa_admin_email,omitempty"`
	Refresh    int    `json:"soa_refresh,omitempty"`
	Retry      int    `json:"soa_retry,omitempty"`
	Expire     int    `json:"soa_expire,omitempty"`
	TTL        int    `json:"soa_ttl,omitempty"`
}

// Map formats a zoneCreateOpts structure into a request body.
func (opts zoneCreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// zoneCreate implements a zone create request.
func zoneCreate(client publicDNSClient, opts optsBuilder) (r zoneCreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getRequestOpts(201)
	resp, err := client.Post(zonesURL(client, zonesAPIPath), &b, &r.Body, reqOpts)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// zoneGet returns information about a zone, given its ID.
func zoneGet(client publicDNSClient, id string) (r zoneGetResult) {
	resp, err := client.Get(getURL(client, zonesAPIPath, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// listOptsBuilder allows extensions to add additional parameters to the
// list request.
type listOptsBuilder interface {
	ToListQuery() (string, error)
}

// listOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned. Marker and Limit are used
// for pagination.
type zoneListOpts struct {
	// Integer value for the limit of values to return.
	Limit int `q:"limit"`

	// UUID of the zone at which you want to set a marker.
	Marker string `q:"marker"`

	ID         string `q:"uuid"`
	Zone       string `q:"zone"`
	PrimaryDNS string `q:"soa_primary_dns"`
	AdminEmail string `q:"soa_admin_email"`
	Serial     int    `q:"soa_serial"`
	Refresh    int    `q:"soa_refresh"`
	Retry      int    `q:"soa_retry"`
	Expire     int    `q:"soa_expire"`
	TTL        int    `q:"soa_ttl"`
	Status     string `q:"status"`
}

// ToListQuery formats a listOpts structure into a query string.
func (opts zoneListOpts) ToListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// zoneList implements a zone list request.
func zoneList(client *gophercloud.ServiceClient, opts listOptsBuilder) (r zoneListResult) {
	url := zonesURL(client, zonesAPIPath)
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
	return
}

// zoneUpdateOpts specifies the attributes to update a zone.
type zoneUpdateOpts struct {
	PrimaryDNS string `json:"soa_primary_dns,omitempty"`
	AdminEmail string `json:"soa_admin_email,omitempty"`
	Refresh    int    `json:"soa_refresh,omitempty"`
	Retry      int    `json:"soa_retry,omitempty"`
	Expire     int    `json:"soa_expire,omitempty"`
	TTL        int    `json:"soa_ttl,omitempty"`
}

// Map formats a zoneUpdateOpts structure into a request body.
func (opts zoneUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// zoneUpdate implements a zone update request.
func zoneUpdate(client publicDNSClient, id string, opts optsBuilder) (r zoneUpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getRequestOpts(200)
	resp, err := client.Put(getURL(client, zonesAPIPath, id), &b, &r.Body, reqOpts)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// zoneDelete implements a zone delete request.
func zoneDelete(client *gophercloud.ServiceClient, id string) (r zoneDeleteResult) {
	resp, err := client.Delete(getURL(client, zonesAPIPath, id), &gophercloud.RequestOpts{
		OkCodes:      []int{204},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
