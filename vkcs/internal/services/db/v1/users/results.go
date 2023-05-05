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
