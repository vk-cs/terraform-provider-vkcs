package serviceusers

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
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

// ServiceUserPage is Pager that is returned from a call to the List function.
type ServiceUserPage struct {
	pagination.PageResult
}

// NextPageURL generates the URL for the page of results after this one.
func (r ServiceUserPage) NextPageURL() (string, error) {
	serviceUsers, err := ExtractServiceUsers(r)
	if err != nil {
		return "", err
	}
	if len(serviceUsers) == 0 {
		return "", nil
	}

	curURL := r.URL
	q := curURL.Query()

	var curOffset int
	if o := q.Get("offset"); o != "" {
		var err error
		curOffset, err = strconv.Atoi(o)
		if err != nil {
			return "", fmt.Errorf("error parsing offset: %w", err)
		}
	}

	offset := curOffset + len(serviceUsers)
	q.Set("offset", strconv.Itoa(offset))
	curURL.RawQuery = q.Encode()

	return curURL.String(), nil
}

// IsEmpty returns true if a ServiceUserPage is empty.
func (r ServiceUserPage) IsEmpty() (bool, error) {
	serviceUsers, err := ExtractServiceUsers(r)
	return len(serviceUsers) == 0, err
}

// GetBody returns the body of a ServiceUserPage.
func (r ServiceUserPage) GetBody() interface{} {
	return r.Body
}

type DeleteResult struct {
	gophercloud.ErrResult
}
