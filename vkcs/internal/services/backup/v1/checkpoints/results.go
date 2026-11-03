package checkpoints

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/plans"
)

type ListCheckpointsGroupedByPlanIDResponse struct {
	ProjectID string               `json:"project_id"`
	PlanID    string               `json:"plan_id"`
	Resources []CheckpointResource `json:"resources"`
}

type CheckpointResource struct {
	ID                         string     `json:"id"`
	Name                       string     `json:"name"`
	ResourceType               string     `json:"resource_type"`
	NumberOfFullBackups        int64      `json:"number_of_full_backups"`
	NumberOfIncrementalBackups int64      `json:"number_of_incremental_backups"`
	LastCreatedAt              string     `json:"last_created_at"`
	TotalSize                  int64      `json:"total_size"`
	Backups                    []CheckpointBackup `json:"backups"`
}

type CheckpointBackup struct {
	CheckpointID         string            `json:"checkpoint_id"`
	Name                 string            `json:"name"`
	Status               string            `json:"status"`
	CheckpointResourceID string            `json:"checkpoint_resource_id"`
	IsIncremental        bool              `json:"is_incremental"`
	IsOldBackup          bool              `json:"is_old_backup"`
	Size                 int64             `json:"size"`
	CreatedAt            string            `json:"created_at"`
	DoneAt               string            `json:"done_at"`
	LinkID               string            `json:"link_id"`
	ObjectLock           *plans.ObjectLock `json:"object_lock,omitempty"`
}

type GetCheckpointResponse struct {
	Checkpoint CheckpointDetails `json:"checkpoint"`
}

type CheckpointDetails struct {
	ID        string                    `json:"id"`
	Resources []CheckpointDetailResource `json:"resources"`
}

type CheckpointDetailResource struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	OriginalID string `json:"original_id"`
	Type       string `json:"type"`
}

type commonResult struct {
	gophercloud.Result
}

func (result commonResult) Extract() (*ListCheckpointsGroupedByPlanIDResponse, error) {
	var response *ListCheckpointsGroupedByPlanIDResponse
	if err := result.ExtractInto(&response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

func (result GetResult) Extract() (*CheckpointDetails, error) {
	var response *GetCheckpointResponse
	if err := result.ExtractInto(&response); err != nil {
		return nil, err
	}
	return &response.Checkpoint, nil
}
