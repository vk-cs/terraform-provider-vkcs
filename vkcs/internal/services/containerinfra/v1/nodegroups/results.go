package nodegroups

import "github.com/gophercloud/gophercloud"

type commonResult struct {
	gophercloud.Result
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type PatchResult struct {
	gophercloud.Result
}

type ResizeResult struct {
	commonResult
}

type ScaleResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}

// Extract parses result into params for node group.
func (r commonResult) Extract() (*NodeGroup, error) {
	var s *NodeGroup
	err := r.ExtractInto(&s)
	return s, err
}

// Extract returns uuid.
func (r PatchResult) Extract() (string, error) {
	var s struct {
		UUID string
	}
	err := r.ExtractInto(&s)
	return s.UUID, err
}

// Extract returns uuid.
func (r ResizeResult) Extract() (string, error) {
	var s struct {
		UUID string `json:"uuid"`
	}
	err := r.ExtractInto(&s)
	return s.UUID, err
}
