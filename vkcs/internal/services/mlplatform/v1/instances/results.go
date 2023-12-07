package instances

import (
	"github.com/gophercloud/gophercloud"
)

type Response struct {
	ID                     string           `json:"id" required:"true"`
	Name                   string           `json:"name" required:"true"`
	Status                 string           `json:"status" required:"true"`
	FlavorID               string           `json:"flavor_id" required:"true"`
	CreatedAt              string           `json:"created_dt" required:"true"`
	PublicIP               string           `json:"public_ip,omitempty"`
	PrivateIP              string           `json:"private_ip,omitempty"`
	DomainName             string           `json:"domain_name,omitempty"`
	InstanceType           string           `json:"instance_type" required:"true"`
	JHAdminName            string           `json:"jh_admin_name" required:"true"`
	MLFlowJHInstanceID     string           `json:"mlflow_jh_instance_id,omitempty"`
	DeployMLFlowInstanceID string           `json:"mlflow_deploy_instance_id,omitempty"`
	Volumes                []VolumeResponse `json:"volumes" required:"true"`
}

type VolumeResponse struct {
	Name             string `json:"name,omitempty"`
	Size             int    `json:"size" required:"true"`
	VolumeType       string `json:"volume_type" required:"true"`
	AvailabilityZone string `json:"availability_zone" required:"true"`
	CinderID         string `json:"cinder_id,omitempty"`
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*Response, error) {
	var res *Response
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

type ActionResult struct {
	commonResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}
