package secrets

import (
	"time"

	"github.com/gophercloud/gophercloud"
)

// ListResult is the result of a list request.
// Call its Extract method to interpret a result as a Secrets object.
type ListResult struct {
	gophercloud.Result
}

type Secrets struct {
	Data struct {
		Keys []string `json:"keys"`
	} `json:"data"`
}

// Extract interprets a get result as a list of secrets' keys.
func (r ListResult) Extract() ([]string, error) {
	var s Secrets
	err := r.ExtractInto(&s)
	return s.Data.Keys, err
}

// GetResult is the result of a get request.
// Call its Extract method to interpret a result as a GetSecret object.
type GetResult struct {
	gophercloud.Result
}

type GetSecret struct {
	Data struct {
		Data     map[string]string `json:"data"`
		Metadata struct {
			CreatedTime time.Time `json:"created_time"`
			Version     int       `json:"version"`
		} `json:"metadata"`
	} `json:"data"`
}

// Extract interprets a get result as a GetSecret object.
func (r GetResult) Extract() (GetSecret, error) {
	var s GetSecret
	err := r.ExtractInto(&s)
	return s, err
}

// GetDeleteProtectionResult is the result of a get request.
// Call its Extract method to interpret a result as a bool.
type GetDeleteProtectionResult struct {
	gophercloud.Result
}

type GetSecretDeleteProtection struct {
	DeleteProtection bool `json:"delete_protection"`
}

// Extract interprets a get result as a bool.
func (r GetDeleteProtectionResult) Extract() (bool, error) {
	var s GetSecretDeleteProtection
	err := r.ExtractInto(&s)
	return s.DeleteProtection, err
}

// CreateOrUpdateResult is the result of a get request.
// Call its Extract method to interpret a result as a CreateOrUpdateSecret object.
type CreateOrUpdateResult struct {
	gophercloud.Result
}

type CreateOrUpdateSecret struct {
	Data struct {
		CreatedTime time.Time `json:"created_time"`
		Version     int       `json:"version"`
	} `json:"data"`
}

// Extract interprets a get result as a CreateOrUpdateSecret object.
func (r CreateOrUpdateResult) Extract() (CreateOrUpdateSecret, error) {
	var s CreateOrUpdateSecret
	err := r.ExtractInto(&s)
	return s, err
}

// SetDeleteProtectionResult is the result of a set delete protection request.
// Call its Extract method to interpret a result as an error.
type SetDeleteProtectionResult struct {
	gophercloud.Result
}

// Extract interprets a get result as an error.
func (r SetDeleteProtectionResult) Extract() error {
	return r.Err
}

// DeleteResult is the result of a delete request.
// Call its Extract method to interpret a result as an error.
type DeleteResult struct {
	gophercloud.Result
}

// Extract interprets a get result as an error.
func (r DeleteResult) Extract() error {
	return r.Err
}
