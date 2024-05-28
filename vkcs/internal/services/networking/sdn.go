package networking

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gophercloud/gophercloud"
)

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

func availableSDNsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("available-sdn")
}

func GetAvailableSDNs(c *gophercloud.ServiceClient) ([]string, error) {
	var sdn []string
	httpResp, err := c.Get(availableSDNsURL(c), &sdn, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting avalible SDN's: %s", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting available SDN's: %s", httpResp.Status)
	}

	for i := 0; i < len(sdn); i++ {
		sdn[i] = strings.ToLower(sdn[i])
	}

	return sdn, nil
}
