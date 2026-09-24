package secrets

import (
	"time"

	"github.com/gophercloud/gophercloud"
)

// GetResult is the result of a get request.
// Call its Extract method to interpret a result as a ListSecrets.
type GetResult struct {
	gophercloud.Result
}

type Secret struct {
	Data struct {
		Data     map[string]string `json:"data"`
		Metadata struct {
			CreatedTime  time.Time  `json:"created_time"`
			DeletionTime *time.Time `json:"deletion_time"`
			Version      int        `json:"version"`
		} `json:"metadata"`
	} `json:"data"`
}

// Extract interprets a get result as a secret.
func (r GetResult) Extract() (Secret, error) {
	var s Secret
	err := r.ExtractInto(&s)
	return s, err
}

// ListResult is the result of a list request.
// Call its Extract method to interpret a result as a Secrets.
type ListResult struct {
	gophercloud.Result
}

type Secrets struct {
	Data struct {
		Keys []string `json:"keys"`
	} `json:"data"`
}

// Extract interprets a get result as a list of secrets.
func (r ListResult) Extract() ([]string, error) {
	var s Secrets
	err := r.ExtractInto(&s)
	return s.Data.Keys, err
}
