package ssldata

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
)

func baseURL(projectID string) string {
	return fmt.Sprintf("projects/%s/sslData", projectID)
}

func sslCertificatesURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(baseURL(projectID))
}

func sslCertificateURL(c *gophercloud.ServiceClient, projectID string, id int) string {
	return c.ServiceURL(baseURL(projectID), strconv.Itoa(id))
}
