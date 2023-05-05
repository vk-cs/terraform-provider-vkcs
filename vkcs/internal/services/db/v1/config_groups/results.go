package configgroups

import (
	"github.com/gophercloud/gophercloud"
)

type ConfigGroupResp struct {
	ID                   string                 `json:"id"`
	DatastoreName        string                 `json:"datastore_name"`
	DatastoreVersionName string                 `json:"datastore_version_name"`
	Name                 string                 `json:"name"`
	Values               map[string]interface{} `json:"values"`
	Updated              string                 `json:"updated"`
	Created              string                 `json:"created"`
	Description          string                 `json:"description"`
}

type ConfigGroupRespOpts struct {
	Configuration *ConfigGroupResp `json:"configuration"`
}

type commonConfigGroupResult struct {
	gophercloud.Result
}

type CreateResult struct {
	commonConfigGroupResult
}

type GetResult struct {
	commonConfigGroupResult
}

type UpdateResult struct {
	gophercloud.ErrResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}

func (r commonConfigGroupResult) Extract() (*ConfigGroupResp, error) {
	var c *ConfigGroupRespOpts
	if err := r.ExtractInto(&c); err != nil {
		return nil, err
	}
	return c.Configuration, nil
}
