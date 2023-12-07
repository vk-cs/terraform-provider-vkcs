package imageinfos

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "image_info"
}

func imageInfosURL(c *gophercloud.ServiceClient, instanceType string) string {
	return fmt.Sprintf("%s?instance_type=%s", c.ServiceURL(baseURL()), instanceType)
}
