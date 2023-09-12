package interfaces

import (
	"github.com/gophercloud/gophercloud"
)

type InterfaceResp struct {
	Interface Interface `json:"dc_interface"`
}

type Interface struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	DCRouterID         string `json:"dc_router_id"`
	NetworkID          string `json:"network_id"`
	SubnetID           string `json:"subnet_id"`
	BGPAnnounceEnabled bool   `json:"bgp_announce_enabled"`
	PortID             string `json:"port_id"`
	SDN                string `json:"sdn"`
	IPAddress          string `json:"ip_address"`
	IPNetmask          int    `json:"ip_netmask"`
	MACAddress         string `json:"mac_address"`
	MTU                int    `json:"mtu"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_ad"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*Interface, error) {
	var res *InterfaceResp
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return &res.Interface, nil
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
