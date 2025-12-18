package s3accounts

import (
	"github.com/gophercloud/gophercloud"
)

func baseURL() string {
	return "s3-accounts"
}

func s3AccountsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func s3AccountURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
