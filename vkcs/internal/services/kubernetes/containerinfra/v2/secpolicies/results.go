package secpolicies

import (
	"github.com/gophercloud/gophercloud"
)

type (
	ClusterSecPolicyResponse struct {
		ClusterSecPolicy ClusterSecPolicy `json:"cluster_security_policy"`
		CreatedAt        string           `json:"created_at"`
		UpdatedAt        string           `json:"updated_at"`
	}

	ClusterSecPolicy struct {
		ID               string                 `json:"id"`
		ClusterID        string                 `json:"cluster_id"`
		SecurityPolicyID string                 `json:"security_policy_id"`
		PolicySettings   string                 `json:"policy_settings"`
		Namespace        string                 `json:"namespace"`
		Enabled          bool                   `json:"enabled"`
		PolicyTemplate   SecurityPolicyTemplate `json:"policy"`
	}

	SecurityPolicyTemplate struct {
		ID                  string `json:"id"`
		Name                string `json:"name"`
		Description         string `json:"description"`
		SettingsDescription string `json:"settings_description"`
		Version             string `json:"version"`
	}

	ClusterSecPolicyID struct {
		ID string `json:"id"`
	}

	CreateResult struct {
		gophercloud.Result
	}

	GetResult struct {
		gophercloud.Result
	}

	UpdateResult struct {
		gophercloud.Result
	}
)

func (r CreateResult) Extract() (string, error) {
	if r.Err != nil {
		return "", r.Err
	}

	var id ClusterSecPolicyID
	err := r.ExtractInto(&id)
	return id.ID, err
}

func (r GetResult) Extract() (*ClusterSecPolicyResponse, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var clusterSecPolicy ClusterSecPolicyResponse
	err := r.ExtractInto(&clusterSecPolicy)
	return &clusterSecPolicy, err
}

func (r UpdateResult) Extract() (*ClusterSecPolicyResponse, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var clusterSecPolicy ClusterSecPolicyResponse
	err := r.ExtractInto(&clusterSecPolicy)
	return &clusterSecPolicy, err
}
