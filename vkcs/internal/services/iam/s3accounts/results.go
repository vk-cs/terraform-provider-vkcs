package s3accounts

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	paginationutil "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/pagination"
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

// ExtractS3Accounts extracts and returnsa S3Account slice.
func ExtractS3Accounts(r pagination.Page) ([]S3Account, error) {
	var s []S3Account
	err := ExtractS3AccountsInto(r, &s)
	return s, err
}

// ExtractS3AccountsInto converts a page into a S3Account slice.
func ExtractS3AccountsInto(r pagination.Page, v any) error {
	return r.(S3AccountPage).ExtractIntoSlicePtr(v, "data")
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

func NewS3AccountPage(r pagination.PageResult) S3AccountPage {
	return S3AccountPage{OffsetPageBase: paginationutil.OffsetPageBase{PageResult: r, Label: "data"}}
}

// S3AccountPage is Pager that is returned from a call to the List function.
type S3AccountPage struct {
	paginationutil.OffsetPageBase
}

type DeleteResult struct {
	gophercloud.ErrResult
}
