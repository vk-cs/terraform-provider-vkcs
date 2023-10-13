package networking

import "github.com/gophercloud/gophercloud"

const (
	NeutronSDN      = "neutron"
	SprutSDN        = "sprut"
	SearchInAllSDNs = "all"
	DefaultSDN      = NeutronSDN
)

type SDNExt struct {
	SDN string `json:"sdn"`
}

func SelectSDN(c *gophercloud.ServiceClient, sdn string) error {
	if sdn != SearchInAllSDNs {
		c.MoreHeaders = map[string]string{
			"X-SDN": sdn,
		}
	}
	return nil
}
