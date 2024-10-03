package origingroups

import "github.com/gophercloud/gophercloud"

type commonResult struct {
	gophercloud.Result
}

// Extract interprets any origin group result as a OriginGroup, if possible.
func (r commonResult) Extract() (*OriginGroup, error) {
	var res OriginGroup
	err := r.ExtractInto(&res)
	return &res, err
}

// CreateResult is the result of a create request. Call its Extract method
// to interpret a result as a OriginGroup.
type CreateResult struct {
	commonResult
}

// GetResult is the result of a get request. Call its Extract method
// to interpret a result as a OriginGroup.
type GetResult struct {
	commonResult
}

// ListResult is the result of a list request. Call its Extract method
// to interpret a result as a slice of OriginGroup.
type ListResult struct {
	gophercloud.Result
}

// Extract interprets a ListResult as a slice of OriginGroup, if possible.
func (r ListResult) Extract() ([]OriginGroup, error) {
	var originGroups []OriginGroup
	err := r.ExtractInto(&originGroups)
	return originGroups, err
}

// UpdateResult is the result of an update request. Call its Extract method
// to interpret a result as a OriginGroup.
type UpdateResult struct {
	commonResult
}

// DeleteResult is the result of a delete request. Call its ExtractErr method
// to determine if a request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// OriginGroup represents a CDN origin group.
type OriginGroup struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	OriginIDs []OriginID `json:"origin_ids"`
	Origins   []Origin   `json:"origins"`
	UseNext   bool       `json:"useNext"`
}

type OriginID struct {
	ID      int    `json:"id"`
	Backup  bool   `json:"backup"`
	Enabled bool   `json:"enabled"`
	Source  string `json:"source"`
}

// Origin represents a CDN origin.
type Origin struct {
	Backup  bool   `json:"backup"`
	Enabled bool   `json:"enabled"`
	Source  string `json:"source"`
}
