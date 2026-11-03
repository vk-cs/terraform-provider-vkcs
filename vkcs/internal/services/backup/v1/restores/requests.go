package restores

import "github.com/gophercloud/gophercloud"

type ListOpts struct {
	Limit  int    `q:"limit"`
	Offset int    `q:"offset"`
	Sort   string `q:"sort"`
}

func (opts ListOpts) ToListQuery() (string, error) {
	query, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return query.String(), nil
}

type RestoreCreate struct {
	Restore *RestoreCreateOpts `json:"restore" required:"true"`
}

type RestoreCreateOpts struct {
	ProviderID   string                       `json:"provider_id" required:"true"`
	CheckpointID string                       `json:"checkpoint_id" required:"true"`
	Parameters   map[string]RestoreParameters `json:"parameters" required:"true"`
}

type RestoreParameters struct {
	RestoreName             string   `json:"restore_name" required:"true"`
	Destination             string   `json:"destination,omitempty"`
	RestoreKeypair          string   `json:"restore_keypair,omitempty"`
	AvailabilityZone        string   `json:"availability_zone,omitempty"`
	RestoreNetID            string   `json:"restore_net_id,omitempty"`
	VolumeType              string   `json:"volume_type,omitempty"`
	VolumeSize              *int64   `json:"volume_size,omitempty"`
	EnableFIP               *bool    `json:"enable_fip,omitempty"`
	Flavor                  string   `json:"flavor,omitempty"`
	VolumeTargetVMID        string   `json:"volume_target_vm_id,omitempty"`
	VolumeAutoresizeEnabled *int64   `json:"volume_autoresize_enabled,omitempty"`
	VolumeAutoresizeMaxSize *int64   `json:"volume_autoresize_max_size,omitempty"`
	SecurityGroups          []string `json:"security_groups,omitempty"`
	CloudMonitoringEnabled  *bool    `json:"cloud_monitoring_enabled,omitempty"`
}

func (opts *RestoreCreate) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

// Create submits a one-time restore request.
func Create(client *gophercloud.ServiceClient, opts *RestoreCreate) (result CreateResult) {
	body, err := opts.Map()
	if err != nil {
		result.Err = err
		return
	}

	response, err := client.Post(restoresURL(client), &body, &result.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, result.Header, result.Err = gophercloud.ParseResponse(response, err)
	return
}

// Get returns restore details and per-resource statuses.
func Get(client *gophercloud.ServiceClient, id string) (result GetResult) {
	response, err := client.Get(restoreURL(client, id), &result.Body, nil)
	_, result.Header, result.Err = gophercloud.ParseResponse(response, err)
	return
}

// List returns one page of restores for the current project.
func List(client *gophercloud.ServiceClient, opts ListOpts) (result ListResult) {
	query, err := opts.ToListQuery()
	if err != nil {
		result.Err = err
		return
	}
	response, err := client.Get(restoresURL(client)+query, &result.Body, nil)
	_, result.Header, result.Err = gophercloud.ParseResponse(response, err)
	return
}

// ListAll returns all restores in the project, using the API's offset pagination.
func ListAll(client *gophercloud.ServiceClient) ([]Restore, error) {
	const pageLimit = 100

	var allRestores []Restore
	for offset := 0; ; offset += pageLimit {
		page, err := List(client, ListOpts{Limit: pageLimit, Offset: offset, Sort: "desc"}).Extract()
		if err != nil {
			return nil, err
		}
		allRestores = append(allRestores, page.Restores...)
		if len(page.Restores) < pageLimit {
			return allRestores, nil
		}
	}
}
