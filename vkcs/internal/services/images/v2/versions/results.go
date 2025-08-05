package versions

import "github.com/gophercloud/gophercloud"

type VersionsResult struct {
	gophercloud.Result
}

func (r VersionsResult) Extract() ([]Version, error) {
	var s VersionsResponse
	err := r.ExtractInto(&s)
	return s.Versions, err
}

type VersionLink struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

type Version struct {
	Status string        `json:"status"`
	ID     string        `json:"id"`
	Links  []VersionLink `json:"links"`
}

type VersionsResponse struct {
	Versions []Version `json:"versions"`
}
