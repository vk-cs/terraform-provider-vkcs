package clusteraddons

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gophercloud/gophercloud"
	v1 "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfraaddons/v1"
)

type ClusterAddon struct {
	ID               string     `json:"cluster_addon_id"`
	Addon            *v1.Addon  `json:"addon"`
	Status           string     `json:"status"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        string     `json:"updated_at"`
	DeletedAt        string     `json:"deleted_at"`
	UserChartValues  string     `json:"user_chart_values"`
	Payload          v1.Payload `json:"payload"`
	UpgradeAvailable bool       `json:"upgrade_available"`
	Migrated         bool       `json:"migrated"`
}

func (a *ClusterAddon) UnmarshalJSON(b []byte) error {
	type tmp ClusterAddon
	var s tmp

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	values := make([]byte, base64.StdEncoding.DecodedLen(len(s.UserChartValues)))
	n, err := base64.StdEncoding.Decode(values, []byte(s.UserChartValues))
	if err != nil {
		return err
	}

	s.UserChartValues = string(values[:n])
	*a = ClusterAddon(s)

	return nil
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*ClusterAddon, error) {
	var s struct {
		Addon *ClusterAddon `json:"addon"`
	}
	err := r.ExtractInto(&s)
	return s.Addon, err
}

type GetResult struct {
	commonResult
}

type UpgradeResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}
