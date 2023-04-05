package vkcs

import (
	"github.com/gophercloud/gophercloud"
)

type commonRecordResult struct {
	gophercloud.Result
	Type string
}

// ExtractA interprets a recordGetResult, recordCreateResult or recordUpdateResult as a recordA.
// An error is returned if the original call or the extraction failed.
func (r commonRecordResult) ExtractA() (*recordA, error) {
	var s recordA
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractAAAA interprets a recordGetResult, recordCreateResult or recordUpdateResult as a recordAAA.
// An error is returned if the original call or the extraction failed.
func (r commonRecordResult) ExtractAAAA() (*recordAAAA, error) {
	var s recordAAAA
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractCNAME interprets a recordGetResult, recordCreateResult or recordUpdateResult as a recordCNAME.
// An error is returned if the original call or the extraction failed.
func (r commonRecordResult) ExtractCNAME() (*recordCNAME, error) {
	var s recordCNAME
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractMX interprets a recordGetResult, recordCreateResult or recordUpdateResult as a recordMX.
// An error is returned if the original call or the extraction failed.
func (r commonRecordResult) ExtractMX() (*recordMX, error) {
	var s recordMX
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractNS interprets a recordGetResult, recordCreateResult or recordUpdateResult as a recordNS.
// An error is returned if the original call or the extraction failed.
func (r commonRecordResult) ExtractNS() (*recordNS, error) {
	var s recordNS
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractSRV interprets a recordGetResult, recordCreateResult or recordUpdateResult as a recordSRV.
// An error is returned if the original call or the extraction failed.
func (r commonRecordResult) ExtractSRV() (*recordSRV, error) {
	var s recordSRV
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractTXT interprets a recordGetResult, recordCreateResult or recordUpdateResult as a recordTXT.
// An error is returned if the original call or the extraction failed.
func (r commonRecordResult) ExtractTXT() (*recordTXT, error) {
	var s recordTXT
	err := r.ExtractInto(&s)
	return &s, err
}

// recordCreateResult is the result of a recordCreate request. Call its Extract method
// to interpret the result as a record.
type recordCreateResult struct {
	commonRecordResult
}

// recordGetResult is the result of a recordGet request. Call its Extract method
// to interpret the result as a record.
type recordGetResult struct {
	commonRecordResult
}

// recordUpdateResult is the result of a recordUpdate request. Call its Extract method
// to interpret the result as a record.
type recordUpdateResult struct {
	commonRecordResult
}

// recordDeleteResult is the result of a recordDelete request. Call its ExtractErr method
// to determine if the request succeeded or failed.
type recordDeleteResult struct {
	gophercloud.ErrResult
}

// recordA represents a public DNS zone record A.
type recordA struct {
	UUID string `json:"uuid"`
	DNS  string `json:"dns"`
	Name string `json:"name"`
	IPv4 string `json:"ipv4"`
	TTL  int    `json:"ttl"`
}

// recordACreateOpts specifies the attributes used to create a record A.
type recordACreateOpts struct {
	Name *string `json:"name" required:"true"`
	IPv4 string  `json:"ipv4" required:"true"`
	TTL  int     `json:"ttl,omitempty"`
}

// recordAAAA represents a public DNS zone record AAAA.
type recordAAAA struct {
	UUID string `json:"uuid"`
	DNS  string `json:"dns"`
	Name string `json:"name"`
	IPv6 string `json:"ipv6"`
	TTL  int    `json:"ttl"`
}

// recordAAAACreateOpts specifies the attributes used to create a record AAAA.
type recordAAAACreateOpts struct {
	Name *string `json:"name" required:"true"`
	IPv6 string  `json:"ipv6" required:"true"`
	TTL  int     `json:"ttl,omitempty"`
}

// recordCNAME represents a public DNS zone record CNAME.
type recordCNAME struct {
	UUID    string `json:"uuid"`
	DNS     string `json:"dns"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

// recordCNAMECreateOpts specifies the attributes used to create a record CNAME.
type recordCNAMECreateOpts struct {
	Name    string `json:"name" required:"true"`
	Content string `json:"content" required:"true"`
	TTL     int    `json:"ttl,omitempty"`
}

// recordMX represents a public DNS zone record MX.
type recordMX struct {
	UUID     string `json:"uuid"`
	DNS      string `json:"dns"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
}

// recordMXCreateOpts specifies the attributes used to create a record MX.
type recordMXCreateOpts struct {
	Name     string `json:"name"`
	Priority int    `json:"priority" required:"true"`
	Content  string `json:"content" required:"true"`
	TTL      int    `json:"ttl,omitempty"`
}

// recordNS represents a public DNS zone record NS.
type recordNS struct {
	UUID    string `json:"uuid"`
	DNS     string `json:"dns"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

// recordNSCreateOpts specifies the attributes used to create a record NS.
type recordNSCreateOpts struct {
	Name    string `json:"name"`
	Content string `json:"content" required:"true"`
	TTL     int    `json:"ttl,omitempty"`
}

// recordSRV represents a public DNS zone record SRV.
type recordSRV struct {
	UUID     string `json:"uuid"`
	DNS      string `json:"dns"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Weight   int    `json:"weight"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	TTL      int    `json:"ttl"`
}

// recordSRVCreateOpts specifies the attributes used to create a record SRV.
type recordSRVCreateOpts struct {
	Name     string `json:"name"`
	Priority int    `json:"priority" required:"true"`
	Weight   int    `json:"weight" required:"true"`
	Host     string `json:"host" required:"true"`
	Port     int    `json:"port" required:"true"`
	TTL      int    `json:"ttl,omitempty"`
}

// recordTXT represents a public DNS zone record TXT.
type recordTXT struct {
	UUID    string `json:"uuid"`
	DNS     string `json:"dns"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

// recordTXTCreateOpts specifies the attributes used to create a record TXT.
type recordTXTCreateOpts struct {
	Name    string `json:"name"`
	Content string `json:"content" required:"true"`
	TTL     int    `json:"ttl,omitempty"`
}

// Map formats a recordACreateOpts structure into a request body.
func (opts recordACreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a recordAAAACreateOpts structure into a request body.
func (opts recordAAAACreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a recordCNAMECreateOpts structure into a request body.
func (opts recordCNAMECreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a recordMXCreateOpts structure into a request body.
func (opts recordMXCreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a recordNSCreateOpts structure into a request body.
func (opts recordNSCreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a recordSRVCreateOpts structure into a request body.
func (opts recordSRVCreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a recordTXTCreateOpts structure into a request body.
func (opts recordTXTCreateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// recordCreate implements a record create request.
func recordCreate(client publicDNSClient, zoneID string, opts optsBuilder, recordType string) (r recordCreateResult) {
	r.Type = recordType
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getRequestOpts(201)
	resp, err := client.Post(recordsURL(client, zonesAPIPath, zoneID, recordType), &b, &r.Body, reqOpts)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// recordGet returns information about a record, given zone ID, its ID, and recordType.
func recordGet(client publicDNSClient, zoneID string, id string, recordType string) (r recordGetResult) {
	r.Type = recordType
	url := recordURL(client, zonesAPIPath, zoneID, recordType, id)
	resp, err := client.Get(url, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// recordAUpdateOpts specifies the attributes used to update a record A.
type recordAUpdateOpts struct {
	Name string `json:"name" required:"true"`
	IPv4 string `json:"ipv4" required:"true"`
	TTL  int    `json:"ttl,omitempty"`
}

// recordAAAAUpdateOpts specifies the attributes used to update a record AAAA.
type recordAAAAUpdateOpts struct {
	Name string `json:"name" required:"true"`
	IPv6 string `json:"ipv6" required:"true"`
	TTL  int    `json:"ttl,omitempty"`
}

// recordCNAMEUpdateOpts specifies the attributes used to update a record CNAME.
type recordCNAMEUpdateOpts struct {
	Name    string `json:"name" required:"true"`
	Content string `json:"content" required:"true"`
	TTL     int    `json:"ttl,omitempty"`
}

// recordMXUpdateOpts specifies the attributes used to update a record MX.
type recordMXUpdateOpts struct {
	Name     string `json:"name" required:"true"`
	Priority int    `json:"priority" required:"true"`
	Content  string `json:"content" required:"true"`
	TTL      int    `json:"ttl,omitempty"`
}

// recordNSUpdateOpts specifies the attributes used to update a record NS.
type recordNSUpdateOpts struct {
	Name    string `json:"name" required:"true"`
	Content string `json:"content" required:"true"`
	TTL     int    `json:"ttl,omitempty"`
}

// recordSRVUpdateOpts specifies the attributes used to update a record SRV.
type recordSRVUpdateOpts struct {
	Name     string `json:"name" required:"true"`
	Priority int    `json:"priority" required:"true"`
	Weight   int    `json:"weight" required:"true"`
	Host     string `json:"host" required:"true"`
	Port     int    `json:"port" required:"true"`
	TTL      int    `json:"ttl,omitempty"`
}

// recordTXTUpdateOpts specifies the attributes used to update a record TXT.
type recordTXTUpdateOpts struct {
	Name    string `json:"name" required:"true"`
	Content string `json:"content" required:"true"`
	TTL     int    `json:"ttl,omitempty"`
}

// Map formats a recordAUpdateOpts structure into a request body.
func (opts recordAUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a recordAAAAUpdateOpts structure into a request body.
func (opts recordAAAAUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a recordCNAMEUpdateOpts structure into a request body.
func (opts recordCNAMEUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a recordMXUpdateOpts structure into a request body.
func (opts recordMXUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a recordNSUpdateOpts structure into a request body.
func (opts recordNSUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a recordSRVUpdateOpts structure into a request body.
func (opts recordSRVUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Map formats a recordTXTUpdateOpts structure into a request body.
func (opts recordTXTUpdateOpts) Map() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// recordUpdate implements a record update request.
func recordUpdate(client publicDNSClient, zoneID string, id string, opts optsBuilder, recordType string) (r recordUpdateResult) {
	r.Type = recordType
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	reqOpts := getRequestOpts(200)
	resp, err := client.Put(recordURL(client, zonesAPIPath, zoneID, recordType, id), &b, &r.Body, reqOpts)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete implements a record delete request.
func recordDelete(client *gophercloud.ServiceClient, zoneID string, id string, recordType string) (r recordDeleteResult) {
	resp, err := client.Delete(recordURL(client, zonesAPIPath, zoneID, recordType, id), &gophercloud.RequestOpts{
		OkCodes:      []int{204},
		JSONResponse: &r.Body,
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
