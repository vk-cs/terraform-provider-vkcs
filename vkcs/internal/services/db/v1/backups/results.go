package backups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

type BackupResp struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	LocationRef string                     `json:"location_ref"`
	InstanceID  string                     `json:"instance_id"`
	ClusterID   string                     `json:"cluster_id"`
	Created     string                     `json:"created"`
	Updated     string                     `json:"updated"`
	Size        float64                    `json:"size"`
	WalSize     float64                    `json:"wal_size"`
	Status      string                     `json:"status"`
	Datastore   *datastores.DatastoreShort `json:"datastore"`
	Meta        string                     `json:"meta"`
}

type BackupRespOpts struct {
	Backup *BackupResp `json:"backup"`
}

type commonResult struct {
	gophercloud.Result
}

type GetResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}

func (r GetResult) Extract() (*BackupResp, error) {
	var b *BackupRespOpts
	if err := r.ExtractInto(&b); err != nil {
		return nil, err
	}
	return b.Backup, nil
}
