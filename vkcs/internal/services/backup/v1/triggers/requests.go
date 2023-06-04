package triggers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type TriggerCreate struct {
	TriggerInfo *CreateOpts `json:"trigger_info" required:"true"`
}

type CreateOpts struct {
	Name       string          `json:"name" required:"true"`
	PlanID     string          `json:"plan_id" required:"true"`
	Properties *PropertiesOpts `json:"properties" required:"true"`
}

type TriggerUpdate struct {
	TriggerInfo *UpdateOpts `json:"trigger_info" required:"true"`
}

type UpdateOpts struct {
	Name       string `json:"name" required:"true"`
	MaxBackups int    `json:"max_backups"`
	Pattern    string `json:"pattern"`
}

// Map converts opts to a map (for a request body)
func (opts *TriggerCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

// Map converts opts to a map (for a request body)
func (opts *TriggerUpdate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

// Create performs request to create backup trigger
func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(triggersURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List will list all backup triggers
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, triggersURL(client),
		func(r pagination.PageResult) pagination.Page {
			return Page{pagination.SinglePageBase(r)}
		})
}

// Update performs request to update backup trigger
func Update(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(triggerURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
