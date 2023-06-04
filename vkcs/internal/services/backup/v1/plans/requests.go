package plans

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// Plan is used to send request to create backup plan
type PlanCreate struct {
	Plan *CreateOpts `json:"plan" required:"true"`
}

// CreateOpts represents parameters of creation of backup plan
type CreateOpts struct {
	Name          string                `json:"name" required:"true"`
	Resources     []*BackupPlanResource `json:"resources" required:"true"`
	ProviderID    string                `json:"provider_id" required:"true"`
	FullDay       *int                  `json:"full_day,omitempty"`
	RetentionType string                `json:"retention_type,omitempty"`
	GFS           *BackupPlanGFS        `json:"gfs,omitempty"`
}

type PlanUpdate struct {
	Plan *UpdateOpts `json:"plan" required:"true"`
}

type UpdateOpts struct {
	Name          string                `json:"name" required:"true"`
	Status        string                `json:"status" required:"true"`
	Resources     []*BackupPlanResource `json:"resources" required:"true"`
	FullDay       int                   `json:"full_day,omitempty"`
	RetentionType string                `json:"retention_type,omitempty"`
	GFS           *BackupPlanGFS        `json:"gfs,omitempty"`
}

// Map converts opts to a map (for a request body)
func (opts *PlanCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

// Map converts opts to a map (for a request body)
func (opts *PlanUpdate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

// Create performs request to create backup plan
func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(plansURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get returns information about backup plan, given its ID
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(planURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Update performs request to update backup plan
func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(planURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return

}

// Delete performs request to delete backup plan
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(planURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
