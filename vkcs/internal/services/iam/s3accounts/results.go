package s3accounts

import (
	"github.com/gophercloud/gophercloud"
)

// CreateResult is the result of a create request. Call its Extract method
// to interpret a result as a CreateS3AccountResponse.
type CreateResult struct {
	gophercloud.Result
}

// Extract interprets a create result as a CreateS3AccountResponse.
func (r CreateResult) Extract() (*CreateS3AccountResponse, error) {
	var s CreateS3AccountResponse
	err := r.ExtractInto(&s)
	return &s, err
}

// GetResult is the result of a get request. Call its Extract method
// to interpret a result as a S3Account.
type GetResult struct {
	gophercloud.Result
}

// Extract interprets a get result as a S3Account.
func (r GetResult) Extract() (*S3Account, error) {
	var s S3Account
	err := r.ExtractInto(&s)
	return &s, err
}

type DeleteResult struct {
	gophercloud.ErrResult
}
