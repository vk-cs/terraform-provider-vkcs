package checkpoints

import (
	"github.com/gophercloud/gophercloud"
)

// PlanCheckpointCreate is the request to create a checkpoint from a plan.
type PlanCheckpointCreate struct {
	PlanCheckpoints *PlanCheckpointCreateOpts `json:"planCheckpoints" required:"true"`
}

// PlanCheckpointCreateOpts contains the plan ID for an on-demand checkpoint.
type PlanCheckpointCreateOpts struct {
	PlanID string `json:"plan_id" required:"true"`
}

// Map converts the plan checkpoint request to a request body.
func (opts *PlanCheckpointCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

// Delete removes a checkpoint.
func Delete(client *gophercloud.ServiceClient, providerID, checkpointID string) (result DeleteResult) {
	response, err := client.Delete(checkpointURL(client, providerID, checkpointID), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, result.Header, result.Err = gophercloud.ParseResponse(response, err)
	return
}

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

// CreateForPlan starts an on-demand checkpoint for a backup plan.
func CreateForPlan(client *gophercloud.ServiceClient, providerID string, opts *PlanCheckpointCreate) (result CreatePlanCheckpointResult) {
	body, err := opts.Map()
	if err != nil {
		result.Err = err
		return
	}

	response, err := client.Post(planCheckpointsURL(client, providerID), &body, &result.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, result.Header, result.Err = gophercloud.ParseResponse(response, err)
	return
}

// Get returns checkpoint details by provider and checkpoint ID.
func Get(client *gophercloud.ServiceClient, providerID, checkpointID string) (result GetResult) {
	response, err := client.Get(checkpointURL(client, providerID, checkpointID), &result.Body, nil)
	_, result.Header, result.Err = gophercloud.ParseResponse(response, err)
	return
}
