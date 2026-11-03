package checkpoints

import "github.com/gophercloud/gophercloud"

type ListOpts struct {
	PlanID string `q:"plan_id"`
}

func (opts ListOpts) ToListQuery() (string, error) {
	query, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return query.String(), nil
}

// ListGroupedByPlanID returns backup checkpoints for one plan, grouped by protected resource.
func ListGroupedByPlanID(client *gophercloud.ServiceClient, opts ListOpts) (result ListResult) {
	query, err := opts.ToListQuery()
	if err != nil {
		result.Err = err
		return
	}

	response, err := client.Get(checkpointsURL(client)+query, &result.Body, nil)
	_, result.Header, result.Err = gophercloud.ParseResponse(response, err)
	return
}

// Get returns checkpoint details by provider and checkpoint ID.
func Get(client *gophercloud.ServiceClient, providerID, checkpointID string) (result GetResult) {
	response, err := client.Get(checkpointURL(client, providerID, checkpointID), &result.Body, nil)
	_, result.Header, result.Err = gophercloud.ParseResponse(response, err)
	return
}
