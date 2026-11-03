package restores

import "github.com/gophercloud/gophercloud"

type RestoreAPIObject struct {
	Restore Restore `json:"restore"`
}

type ListRestoreResponse struct {
	Restores []Restore `json:"restores"`
}

type Restore struct {
	ID              string                     `json:"id"`
	ProviderID      string                     `json:"provider_id"`
	CheckpointID    string                     `json:"checkpoint_id"`
	Parameters      map[string]RestoreParameters `json:"parameters"`
	ResourcesStatus map[string]string          `json:"resources_status,omitempty"`
	ResourcesReason map[string]string          `json:"resources_reason,omitempty"`
}

type commonResult struct {
	gophercloud.Result
}

func (result commonResult) Extract() (*Restore, error) {
	var response *RestoreAPIObject
	if err := result.ExtractInto(&response); err != nil {
		return nil, err
	}
	return &response.Restore, nil
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type ListResult struct {
	commonResult
}

func (result ListResult) Extract() (*ListRestoreResponse, error) {
	var response *ListRestoreResponse
	if err := result.ExtractInto(&response); err != nil {
		return nil, err
	}
	return response, nil
}
