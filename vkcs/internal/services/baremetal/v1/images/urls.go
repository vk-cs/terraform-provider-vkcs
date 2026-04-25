package images

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "images"
}

func imagesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func imageURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
