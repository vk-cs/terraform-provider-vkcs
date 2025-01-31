package anycastips

import "github.com/gophercloud/gophercloud"

type commonResult struct {
	gophercloud.Result
}

// Extract interprets any anycast IP result as an AnycastIP, if possible.
func (r commonResult) Extract() (*AnycastIP, error) {
	var res AnycastIPResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.AnycastIP, nil
}

// CreateResult is the result of a create request. Call its Extract method
// to interpret a result as an AnycastIP.
type CreateResult struct {
	commonResult
}

// GetResult is the result of a get request. Call its Extract method
// to interpret a result as an AnycastIP.
type GetResult struct {
	commonResult
}

// UpdateResult is the result of an update request. Call its Extract method
// to interpret a result as an AnycastIP.
type UpdateResult struct {
	commonResult
}

// DeleteResult is the result of a delete request. Call its ExtractErr method
// to determine if a request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

type AnycastIPResp struct {
	AnycastIP AnycastIP `json:"anycastip"`
}

// AnycastIP represents an anycast IP.
type AnycastIP struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	NetworkID    string                 `json:"network_id"`
	SubnetID     string                 `json:"subnet_id"`
	IPAddress    string                 `json:"ip_address"`
	Associations []AnycastIPAssociation `json:"associations,omitempty"`
	HealthCheck  *AnycastIPHealthCheck  `json:"healthcheck,omitempty"`
}
