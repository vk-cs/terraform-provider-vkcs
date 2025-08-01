package versions

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
