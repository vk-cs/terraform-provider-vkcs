package plans

import "github.com/gophercloud/gophercloud"

type PlanResp struct {
	Plan PlanResponse `json:"plan"`
}

type PlanResponse struct {
	ID            string                `json:"id" required:"true"`
	Name          string                `json:"name" required:"true"`
	Resources     []*BackupPlanResource `json:"resources" required:"true"`
	ProviderID    string                `json:"provider_id" required:"true"`
	Status        string                `json:"status" required:"true"`
	FullDay       *int                  `json:"full_day,omitempty"`
	RetentionType string                `json:"retention_type,omitempty"`
	GFS           *BackupPlanGFS        `json:"gfs,omitempty"`
}

// BackupPlanResource represents a backup plan resource info
type BackupPlanResource struct {
	ID   string `json:"id" required:"true"`
	Type string `json:"type" required:"true"`
	Name string `json:"name" required:"true"`
}

// BackupPlanGFS represents a backup plan gfs policy info
type BackupPlanGFS struct {
	Grandfather int `json:"grandfather,omitempty"`
	Father      int `json:"father,omitempty"`
	Son         int `json:"son" required:"true"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*PlanResponse, error) {
	var s *PlanResp
	err := r.ExtractInto(&s)
	return &s.Plan, err
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}
