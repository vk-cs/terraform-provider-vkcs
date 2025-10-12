package securitypolicytemplates

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "security_policy"
}

func securityPolicyTemplatesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}
