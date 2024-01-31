package securitypolicies

import (
	"github.com/gophercloud/gophercloud"
)

type SecurityPolicy struct {
	UUID                     string `json:"uuid" required:"true"`
	ClusterID                string `json:"cluster_uuid" required:"true"`
	SecurityPolicyTemplateID string `json:"security_policy_uuid" required:"true"`
	PolicySettings           string `json:"policy_settings"`
	Namespace                string `json:"namespace"`
	Enabled                  bool   `json:"enabled"`
	CreatedAt                string `json:"created_at"`
	UpdatedAt                string `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

// CreateResult is the response of a Create operations.
type CreateResult struct {
	commonResult
}

// DeleteResult is the result from a Delete operation. Call its Extract or ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// UpdateResult is the response of a Update operations.
type UpdateResult struct {
	commonResult
}

// Extract parses result into params for security policy.
func (r commonResult) Extract() (*SecurityPolicy, error) {
	var s *SecurityPolicy
	err := r.ExtractInto(&s)
	return s, err
}
