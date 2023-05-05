package databases

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Database represents a Database API resource.
type Database struct {
	// Specifies the name of the MySQL database.
	Name string

	// Set of symbols and encodings. The default character set is utf8.
	CharSet string

	// Set of rules for comparing characters in a character set. The default
	// value for collate is utf8_general_ci.
	Collate string
}

// Custom type implementation of DB Page
type Page struct {
	pagination.LinkedPageBase
}

// IsEmpty checks to see whether the collection is empty.
func (page Page) IsEmpty() (bool, error) {
	dbs, err := ExtractDatabases(page)
	return len(dbs) == 0, err
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

// ExtractDatabases will convert a generic pagination struct into a more
// relevant slice of DB structs.
func ExtractDatabases(page pagination.Page) ([]Database, error) {
	r := page.(Page)
	var s struct {
		Databases []Database `json:"databases"`
	}
	err := r.ExtractInto(&s)
	return s.Databases, err
}

type commonDatabaseResult struct {
	gophercloud.ErrResult
}

// CreateResult represents result of database create
type CreateResult struct {
	commonDatabaseResult
}

// DeleteResult represents result of database delete
type DeleteResult struct {
	commonDatabaseResult
}
