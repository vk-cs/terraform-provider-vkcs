package users

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/databases"
)

// User represents a database user
type User struct {
	// The user name
	Name string

	// The user password
	Password string

	// The databases associated with this user
	Databases []databases.Database
}

// UserRespOpts is used to get user response
type UserRespOpts struct {
	User *User `json:"user"`
}

// Custom type implementation of gophercloud/users.UserPage
type Page struct {
	pagination.LinkedPageBase
}

// IsEmpty checks to see whether the collection is empty.
func (page Page) IsEmpty() (bool, error) {
	users, err := ExtractUsers(page)
	return len(users) == 0, err
}

// NextPageURL will retrieve the next page URL.
func (page Page) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractUsers will convert a generic pagination struct into a more
// relevant slice of User structs.
func ExtractUsers(r pagination.Page) ([]User, error) {
	var s struct {
		Users []User `json:"users"`
	}
	err := (r.(Page)).ExtractInto(&s)
	return s.Users, err
}

type commonResult struct {
	gophercloud.ErrResult
}

type commonUserResult struct {
	gophercloud.Result
}

// GetResult represents result of database user get
type GetResult struct {
	commonUserResult
}

// CreateResult represents result of database user create
type CreateResult struct {
	commonResult
}

// UpdateResult represents result of database user update
type UpdateResult struct {
	commonResult
}

// UpdateDatabasesResult represents result of database user database update
type UpdateDatabasesResult struct {
	commonResult
}

// DeleteDatabaseResult represents result of database user delete
type DeleteDatabaseResult struct {
	commonResult
}

// DeleteResult represents result of database user delete
type DeleteResult struct {
	commonResult
}

// Extract is used to extract result into user response struct.
func (r GetResult) Extract() (*User, error) {
	var u *UserRespOpts
	if err := r.ExtractInto(&u); err != nil {
		return nil, err
	}
	if u != nil && u.User != nil {
		return u.User, nil
	}

	// Fallback for APIs that return users array even for GET.
	var list struct {
		Users []User `json:"users"`
	}
	if err := r.ExtractInto(&list); err != nil {
		return nil, err
	}
	if len(list.Users) > 0 {
		return &list.Users[0], nil
	}
	return nil, nil
}
