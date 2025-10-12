package securitypolicies

import "github.com/gophercloud/gophercloud"

func baseURL() string {
	return "cluster_security_policy"
}

func securityPoliciesURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(baseURL())
}

func securityPolicyURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(baseURL(), id)
}
