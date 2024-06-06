package ipportforwardings

import (
	"github.com/gophercloud/gophercloud"
)

type IPPortForwardingResp struct {
	IPPortForwarding IPPortForwarding `json:"dc_ip_port_forwarding"`
}

type IPPortForwarding struct {
	DCInterfaceID string  `json:"dc_interface_id"`
	Protocol      string  `json:"protocol"`
	Source        *string `json:"source"`
	Destination   *string `json:"destination"`
	Port          *int64  `json:"port"`
	ToDestination string  `json:"to_destination"`
	ToPort        *int64  `json:"to_port"`
	Name          string  `json:"name,omitempty"`
	Description   string  `json:"description,omitempty"`
	CreatedAt     string  `json:"created_at"`
	ID            string  `json:"id"`
	UpdatedAt     string  `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*IPPortForwarding, error) {
	var res *IPPortForwardingResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.IPPortForwarding, nil
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}
