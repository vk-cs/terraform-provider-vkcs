package instances

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

type CreateOpts struct {
	InstanceName           string             `json:"instance_name" required:"true"`
	DomainName             string             `json:"domain_name"`
	InstanceType           string             `json:"instance_type" required:"true"`
	JHAdminName            string             `json:"jh_admin_name,omitempty"`
	JHAdminPassword        string             `json:"jh_admin_password,omitempty"`
	MLFlowJHInstanceID     string             `json:"mlflow_jh_instance_id,omitempty"`
	ISMLFlowDemoMode       bool               `json:"is_mlflow_demo_mode,omitempty"`
	DeployMLFlowInstanceID string             `json:"deploy_mlflow_instance_id,omitempty"`
	Flavor                 string             `json:"flavor" required:"true"`
	Volumes                []VolumeCreateOpts `json:"volumes" required:"true"`
	Networks               NetworkCreateOpts  `json:"networks" required:"true"`
	S3FSBucket             string             `json:"s3fs_bucket,omitempty"`
	ISGPU                  bool               `json:"is_gpu"`
}

type VolumeCreateOpts struct {
	Name             string `json:"name,omitempty"`
	Size             int    `json:"size" required:"true"`
	VolumeType       string `json:"volume_type" required:"true"`
	AvailabilityZone string `json:"availability_zone" required:"true"`
}

type NetworkCreateOpts struct {
	IPPool    string `json:"ip_pool,omitempty"`
	NetworkID string `json:"network_id" required:"true"`
}

type ActionOpts struct {
	Action ResizeAction `json:"action" required:"true"`
}

type ResizeAction struct {
	Resize ResizeActionOpts `json:"resize" required:"true"`
}

type ResizeActionOpts struct {
	Type   string             `json:"type" required:"true"`
	Params ResizeActionParams `json:"params" required:"true"`
}

type ResizeActionParams struct {
	Flavor  string               `json:"flavor,omitempty"`
	Volumes []ResizeVolumeParams `json:"volumes,omitempty"`
}

type ResizeVolumeParams struct {
	ID   string `json:"id" required:"true"`
	Size int    `json:"size" required:"true"`
}

func (opts *CreateOpts) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func (opts *ActionOpts) Map() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(*opts, "")
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(instancesURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(instanceURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Action(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r ActionResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(instanceActionURL(client, id), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(instanceURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
