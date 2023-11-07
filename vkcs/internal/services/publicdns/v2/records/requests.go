package records

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

// RecordA represents a public DNS zone record A.
type RecordA struct {
	UUID string `json:"uuid"`
	DNS  string `json:"dns"`
	Name string `json:"name"`
	IPv4 string `json:"ipv4"`
	TTL  int    `json:"ttl"`
}

// RecordAAAA represents a public DNS zone record AAAA.
type RecordAAAA struct {
	UUID string `json:"uuid"`
	DNS  string `json:"dns"`
	Name string `json:"name"`
	IPv6 string `json:"ipv6"`
	TTL  int    `json:"ttl"`
}

// RecordCNAME represents a public DNS zone record CNAME.
type RecordCNAME struct {
	UUID    string `json:"uuid"`
	DNS     string `json:"dns"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

// RecordMX represents a public DNS zone record MX.
type RecordMX struct {
	UUID     string `json:"uuid"`
	DNS      string `json:"dns"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
}

// RecordNS represents a public DNS zone record NS.
type RecordNS struct {
	UUID    string `json:"uuid"`
	DNS     string `json:"dns"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

// RecordSRV represents a public DNS zone record SRV.
type RecordSRV struct {
	UUID     string `json:"uuid"`
	DNS      string `json:"dns"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Weight   int    `json:"weight"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	TTL      int    `json:"ttl"`
}

// RecordTXT represents a public DNS zone record TXT.
type RecordTXT struct {
	UUID    string `json:"uuid"`
	DNS     string `json:"dns"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type CreateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// RecordACreateOpts specifies the attributes used to create a record A.
type RecordACreateOpts struct {
	Name *string `json:"name" required:"true"`
	IPv4 string  `json:"ipv4" required:"true"`
	TTL  int     `json:"ttl,omitempty"`
}

// RecordAAAACreateOpts specifies the attributes used to create a record AAAA.
type RecordAAAACreateOpts struct {
	Name *string `json:"name" required:"true"`
	IPv6 string  `json:"ipv6" required:"true"`
	TTL  int     `json:"ttl,omitempty"`
}

// RecordCNAMECreateOpts specifies the attributes used to create a record CNAME.
type RecordCNAMECreateOpts struct {
	Name    string `json:"name" required:"true"`
	Content string `json:"content" required:"true"`
	TTL     int    `json:"ttl,omitempty"`
}

// RecordMXCreateOpts specifies the attributes used to create a record MX.
type RecordMXCreateOpts struct {
	Name     string `json:"name"`
	Priority int    `json:"priority" required:"true"`
	Content  string `json:"content" required:"true"`
	TTL      int    `json:"ttl,omitempty"`
}

// RecordNSCreateOpts specifies the attributes used to create a record NS.
type RecordNSCreateOpts struct {
	Name    string `json:"name"`
	Content string `json:"content" required:"true"`
	TTL     int    `json:"ttl,omitempty"`
}

// RecordSRVCreateOpts specifies the attributes used to create a record SRV.
type RecordSRVCreateOpts struct {
	Name     string `json:"name"`
	Priority int    `json:"priority" required:"true"`
	Weight   int    `json:"weight" required:"true"`
	Host     string `json:"host" required:"true"`
	Port     int    `json:"port" required:"true"`
	TTL      int    `json:"ttl,omitempty"`
}

// RecordTXTCreateOpts specifies the attributes used to create a record TXT.
type RecordTXTCreateOpts struct {
	Name    string `json:"name"`
	Content string `json:"content" required:"true"`
	TTL     int    `json:"ttl,omitempty"`
}

// Map formats a RecordACreateOpts structure into a request body.
func (opts RecordACreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a RecordAAAACreateOpts structure into a request body.
func (opts RecordAAAACreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a RecordCNAMECreateOpts structure into a request body.
func (opts RecordCNAMECreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a RecordMXCreateOpts structure into a request body.
func (opts RecordMXCreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a RecordNSCreateOpts structure into a request body.
func (opts RecordNSCreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a RecordSRVCreateOpts structure into a request body.
func (opts RecordSRVCreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a RecordTXTCreateOpts structure into a request body.
func (opts RecordTXTCreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Create implements a record create request.
func Create(client *gophercloud.ServiceClient, zoneID string, opts CreateOptsBuilder, recordType string) (r CreateResult) {
	r.Type = recordType
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(recordsURL(client, zoneID, recordType), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}

// Get returns information about a record, given zone ID, its ID, and recordType.
func Get(client *gophercloud.ServiceClient, zoneID string, id string, recordType string) (r GetResult) {
	r.Type = recordType
	url := recordURL(client, zoneID, recordType, id)
	resp, err := client.Get(url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}

type UpdateOptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// RecordAUpdateOpts specifies the attributes used to update a record A.
type RecordAUpdateOpts struct {
	Name string `json:"name" required:"true"`
	IPv4 string `json:"ipv4" required:"true"`
	TTL  int    `json:"ttl,omitempty"`
}

// RecordAAAAUpdateOpts specifies the attributes used to update a record AAAA.
type RecordAAAAUpdateOpts struct {
	Name string `json:"name" required:"true"`
	IPv6 string `json:"ipv6" required:"true"`
	TTL  int    `json:"ttl,omitempty"`
}

// RecordCNAMEUpdateOpts specifies the attributes used to update a record CNAME.
type RecordCNAMEUpdateOpts struct {
	Name    string `json:"name" required:"true"`
	Content string `json:"content" required:"true"`
	TTL     int    `json:"ttl,omitempty"`
}

// RecordMXUpdateOpts specifies the attributes used to update a record MX.
type RecordMXUpdateOpts struct {
	Name     string `json:"name" required:"true"`
	Priority int    `json:"priority" required:"true"`
	Content  string `json:"content" required:"true"`
	TTL      int    `json:"ttl,omitempty"`
}

// RecordNSUpdateOpts specifies the attributes used to update a record NS.
type RecordNSUpdateOpts struct {
	Name    string `json:"name" required:"true"`
	Content string `json:"content" required:"true"`
	TTL     int    `json:"ttl,omitempty"`
}

// RecordSRVUpdateOpts specifies the attributes used to update a record SRV.
type RecordSRVUpdateOpts struct {
	Name     string `json:"name" required:"true"`
	Priority int    `json:"priority" required:"true"`
	Weight   int    `json:"weight" required:"true"`
	Host     string `json:"host" required:"true"`
	Port     int    `json:"port" required:"true"`
	TTL      int    `json:"ttl,omitempty"`
}

// RecordTXTUpdateOpts specifies the attributes used to update a record TXT.
type RecordTXTUpdateOpts struct {
	Name    string `json:"name" required:"true"`
	Content string `json:"content" required:"true"`
	TTL     int    `json:"ttl,omitempty"`
}

// Map formats a RecordAUpdateOpts structure into a request body.
func (opts RecordAUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a RecordAAAAUpdateOpts structure into a request body.
func (opts RecordAAAAUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a RecordCNAMEUpdateOpts structure into a request body.
func (opts RecordCNAMEUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a RecordMXUpdateOpts structure into a request body.
func (opts RecordMXUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a RecordNSUpdateOpts structure into a request body.
func (opts RecordNSUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a RecordSRVUpdateOpts structure into a request body.
func (opts RecordSRVUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a RecordTXTUpdateOpts structure into a request body.
func (opts RecordTXTUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Update implements a record update request.
func Update(client *gophercloud.ServiceClient, zoneID string, id string, opts UpdateOptsBuilder, recordType string) (r UpdateResult) {
	r.Type = recordType
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(recordURL(client, zoneID, recordType, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}

// Delete implements a record delete request.
func Delete(client *gophercloud.ServiceClient, zoneID string, id string, recordType string) (r DeleteResult) {
	resp, err := client.Delete(recordURL(client, zoneID, recordType, id), &gophercloud.RequestOpts{
		OkCodes:      []int{204},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = util.ErrorWithRequestID(r.Err, r.Header.Get(util.RequestIDHeader))
	return
}
