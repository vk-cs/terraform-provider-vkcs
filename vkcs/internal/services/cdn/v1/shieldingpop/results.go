package shieldingpop

import "github.com/gophercloud/gophercloud"

// ListResult is th result of a list request. Call its Extract method
// to interpret the result as a slice of ShieldingPop.
type ListResult struct {
	gophercloud.Result
}

// Extract interprets a ListResult as a slice of ShieldingPop, if possible.
func (r ListResult) Extract() ([]ShieldingPop, error) {
	var shieldingPops []ShieldingPop
	err := r.ExtractInto(&shieldingPops)
	return shieldingPops, err
}

// ShieldingPop represents a CDN origin shielding point of presence.
type ShieldingPop struct {
	ID         int    `json:"id"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Datacenter string `json:"datacenter"`
}
