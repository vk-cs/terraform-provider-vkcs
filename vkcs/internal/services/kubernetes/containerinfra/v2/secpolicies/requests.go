package secpolicies

import (
	"github.com/gophercloud/gophercloud"
)

// CreateOpts represents options for creating a security policy in cluster
type CreateOpts struct {
	ClusterID        string
	SecurityPolicyID string `json:"security_policy_id" required:"true"`
	PolicySettings   string `json:"policy_settings" required:"true"`
	Namespace        string `json:"namespace" required:"true"`
	Enabled          bool   `json:"enabled" required:"true"`
}

// UpdateOpts
type UpdateOpts struct {
	ClusterID               string
	ClusterSecurityPolicyID string
	PolicySettings          string `json:"policy_settings" required:"true"`
	Namespace               string `json:"namespace" required:"true"`
	Enabled                 bool   `json:"enabled" required:"true"`
}

func (opts CreateOpts) ToSecurityPolicyCreateMap() (map[string]interface{}, error) {
	return map[string]interface{}{
		"security_policy_id": opts.SecurityPolicyID,
		"policy_settings":    opts.PolicySettings,
		"namespace":          opts.Namespace,
		"enabled":            opts.Enabled,
	}, nil
}

func (opts UpdateOpts) ToSecurityPolicyUpdateMap() (map[string]interface{}, error) {
	return map[string]interface{}{
		"policy_settings": opts.PolicySettings,
		"namespace":       opts.Namespace,
		"enabled":         opts.Enabled,
	}, nil
}

// Create creates a new security policy in the specified cluster
func Create(client *gophercloud.ServiceClient, opts CreateOpts) (res CreateResult) {
	reqBody, err := opts.ToSecurityPolicyCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Post(rootURL(client, opts.ClusterID), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

// Get retrieves a specific cluster sec policy group on its ID
func Get(client *gophercloud.ServiceClient, clusterSecPolicyID string) (res GetResult) {
	_, res.Err = client.Get(getURL(client, clusterSecPolicyID), &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

// Update updates a specific cluster sec policy with new configuration
func Update(client *gophercloud.ServiceClient, opts UpdateOpts) (res UpdateResult) {
	reqBody, err := opts.ToSecurityPolicyUpdateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Patch(resourceURL(client, opts.ClusterID, opts.ClusterSecurityPolicyID), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}

// Delete deletes a specific cluster sec policy group
func Delete(client *gophercloud.ServiceClient, clusterID, clusterSecPolicyID string) error {
	_, err := client.Delete(resourceURL(client, clusterID, clusterSecPolicyID), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return err
}
