package clusters

import "github.com/gophercloud/gophercloud"

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

// Extract is a function that accepts a result and extracts a cluster resource.
func (r commonResult) Extract() (*Cluster, error) {
	var s *Cluster
	err := r.ExtractInto(&s)
	return s, err
}

// UpdateResult is the response of a Update operations.
type UpdateResult struct {
	commonResult
}

// UpgradeResult is the response of a Upgrade operations.
type UpgradeResult struct {
	commonResult
}

// ResizeResult is the response of a Resize operations.
type ResizeResult struct {
	commonResult
}

func (r CreateResult) Extract() (string, error) {
	var s struct {
		UUID string
	}
	err := r.ExtractInto(&s)
	return s.UUID, err
}

func (r UpdateResult) Extract() (string, error) {
	var s struct {
		UUID string
	}
	err := r.ExtractInto(&s)
	return s.UUID, err
}
