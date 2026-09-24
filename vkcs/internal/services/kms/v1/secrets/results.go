package secrets

import "github.com/gophercloud/gophercloud"

// ListResult is the result of a get request.
// Call its Extract method to interpret a result as a ListSecrets.
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
