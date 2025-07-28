package products

import (
	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "products"
}

func productsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}
