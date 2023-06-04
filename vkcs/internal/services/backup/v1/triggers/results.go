package triggers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type TriggerResp struct {
	TriggerInfo TriggerResponse `json:"trigger_info"`
}

type TriggerResponse struct {
	ID         string          `json:"id" required:"true"`
	PlanID     string          `json:"plan_id" required:"true"`
	Name       string          `json:"name" required:"true"`
	Properties *PropertiesOpts `json:"properties" required:"true"`
}

type PropertiesOpts struct {
	Pattern    string `json:"pattern" required:"true"`
	MaxBackups int    `json:"max_backups,omitempty"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*TriggerResponse, error) {
	var s *TriggerResp
	err := r.ExtractInto(&s)
	return &s.TriggerInfo, err
}

type CreateResult struct {
	commonResult
}

// Page represents a page of backup triggers
type Page struct {
	pagination.SinglePageBase
}

// IsEmpty indicates whether a backup trigger collection is empty.
func (r Page) IsEmpty() (bool, error) {
	tr, err := ExtractTriggers(r)
	return len(tr) == 0, err
}

// ExtractTriggers retrieves a slice of backup TriggerResponse structs from a paginated
// collection.
func ExtractTriggers(r pagination.Page) ([]TriggerResponse, error) {
	var s struct {
		Triggers []TriggerResponse `json:"triggers"`
	}
	err := (r.(Page)).ExtractInto(&s)
	return s.Triggers, err
}
