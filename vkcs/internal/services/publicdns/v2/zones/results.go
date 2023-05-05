package zones

import "github.com/gophercloud/gophercloud"

type commonResult struct {
	gophercloud.Result
}

// Extract interprets a zoneGetResult, zoneCreateResult or zoneUpdateResult as a zone.
// An error is returned if the original call or the extraction failed.
func (r commonResult) Extract() (*Zone, error) {
	var s *Zone
	err := r.ExtractInto(&s)
	return s, err
}

// zoneCreateResult is the result of a zoneCreate request. Call its Extract method
// to interpret the result as a zone.
type CreateResult struct {
	commonResult
}

// GetResult is the result of a Get request. Call its Extract method
// to interpret the result as a Zone.
type GetResult struct {
	commonResult
}

// ListResult is the result of a List request. Call its Extract method
// to interpret the result as a slice of Zones.
type ListResult struct {
	gophercloud.Result
}

// Extract extracts a slice of zones from a zoneListResult.
func (r ListResult) Extract() ([]Zone, error) {
	var zones []Zone
	err := r.ExtractInto(&zones)
	return zones, err
}

// zoneUpdateResult is the result of a zoneUpdate request. Call its Extract method
// to interpret the result as a zone.
type UpdateResult struct {
	commonResult
}

// zoneDeleteResult is the result of a zoneDelete request. Call its ExtractErr method
// to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Zone represents a public DNS zone.
type Zone struct {
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
