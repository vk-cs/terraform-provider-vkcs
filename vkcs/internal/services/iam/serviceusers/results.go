package serviceusers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	paginationutil "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/pagination"
)

// CreateResult is the result of a create request. Call its Extract method
// to interpret a result as a CreateServiceUserResponse.
type CreateResult struct {
	gophercloud.Result
}

// Extract interprets a create result as a CreateServiceUserResponse.
func (r CreateResult) Extract() (*CreateServiceUserResponse, error) {
	var s CreateServiceUserResponse
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractServiceUsers extracts and returns ServiceUsers.
func ExtractServiceUsers(r pagination.Page) ([]ServiceUser, error) {
	var s []ServiceUser
	err := ExtractServiceUsersInto(r, &s)
	return s, err
}

// ExtractServiceUsersInto converts a page into a slice of ServiceUsers.
func ExtractServiceUsersInto(r pagination.Page, v interface{}) error {
	return r.(ServiceUserPage).Result.ExtractIntoSlicePtr(v, "data")
}

// GetResult is the result of a get request. Call its Extract method
// to interpret a result as a ServiceUser.
type GetResult struct {
	gophercloud.Result
}

// Extract interprets a get result as a ServiceUser.
func (r GetResult) Extract() (*ServiceUser, error) {
	var s ServiceUser
	err := r.ExtractInto(&s)
	return &s, err
}

func NewServiceUserPage(r pagination.PageResult) ServiceUserPage {
	return ServiceUserPage{OffsetPageBase: paginationutil.OffsetPageBase{PageResult: r, Label: "data"}}
}

// ServiceUserPage is Pager that is returned from a call to the List function.
type ServiceUserPage struct {
	paginationutil.OffsetPageBase
}

type DeleteResult struct {
	gophercloud.ErrResult
}
