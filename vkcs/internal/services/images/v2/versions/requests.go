package versions

import "github.com/gophercloud/gophercloud"

func Get(client *gophercloud.ServiceClient) ([]Version, error) {
	var result VersionsResponse

	_, err := client.Get(versionsURL(client), &result, &gophercloud.RequestOpts{
		OkCodes: []int{200, 300},
	})
	if err != nil {
		return nil, err
	}

	return result.Versions, nil
}
