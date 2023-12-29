package backups

import (
	"github.com/gophercloud/gophercloud"
)

type Response struct {
	ID         string `json:"id" required:"true"`
	InstanceID string `json:"instance_id" required:"true"`
	CinderID   string `json:"cinder_id" required:"true"`
	Status     string `json:"status" required:"true"`
	CreatedAt  string `json:"created_at" required:"true"`
	BackupID   string `json:"backup_id" required:"true"`
	Comment    string `json:"comment,omitempty"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() ([]*Response, error) {
	var res []*Response
	if err := r.ExtractInto(&res); err != nil {
		return nil, err
	}
	return res, nil
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}
