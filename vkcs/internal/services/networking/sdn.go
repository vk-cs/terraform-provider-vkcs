package networking

import "github.com/gophercloud/gophercloud"

const (
	NeutronSDN      = "neutron"
	SprutSDN        = "sprut"
	SearchInAllSDNs = "all"
	DefaultSDN      = NeutronSDN
)

func SelectSDN(c *gophercloud.ServiceClient, sdn string) error {
	if sdn != SearchInAllSDNs {
		c.MoreHeaders = map[string]string{
			"X-SDN": sdn,
		}
	}
	return nil
}
