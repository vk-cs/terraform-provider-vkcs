package records

import (
	"strings"

	"github.com/gophercloud/gophercloud"
)

func baseURL(zoneID string, recordType string) string {
	parts := []string{"dns", zoneID, strings.ToLower(recordType)}
	return strings.Join(parts, "/")
}

func recordsURL(c *gophercloud.ServiceClient, zoneID string, recordType string) string {
	return c.ServiceURL(baseURL(zoneID, recordType), "")
}

func recordURL(c *gophercloud.ServiceClient, zoneID string, recordType string, recordID string) string {
	return c.ServiceURL(baseURL(zoneID, recordType), recordID)
}
