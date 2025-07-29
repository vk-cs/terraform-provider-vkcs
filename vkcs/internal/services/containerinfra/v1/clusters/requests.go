package clusters

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type OptsBuilder interface {
	Map() (map[string]interface{}, error)
}

// CreateOpts contains options to create cluster
type CreateOpts struct {
	ClusterTemplateID         string            `json:"cluster_template_id" required:"true"`
	Keypair                   string            `json:"keypair,omitempty"`
	Labels                    map[string]string `json:"labels,omitempty"`
	MasterCount               int               `json:"master_count,omitempty"`
	MasterFlavorID            string            `json:"master_flavor_id,omitempty"`
	Name                      string            `json:"name"`
	NetworkID                 string            `json:"network_id" required:"true"`
	SubnetID                  string            `json:"subnet_id" required:"true"`
	PodsNetworkCidr           string            `json:"pods_network_cidr,omitempty"`
	FloatingIPEnabled         bool              `json:"floating_ip_enabled"`
	APILBFIP                  string            `json:"api_lb_fip,omitempty"`
	RegistryAuthPassword      string            `json:"registry_auth_password,omitempty"`
	AvailabilityZone          string            `json:"availability_zone,omitempty"`
	AvailabilityZones         []string          `json:"availability_zones,omitempty"`
	LoadbalancerSubnetID      string            `json:"loadbalancer_subnet_id,omitempty"`
	InsecureRegistries        []string          `json:"insecure_registries,omitempty"`
	DNSDomain                 string            `json:"dns_domain,omitempty"`
	SecurityPolicySyncEnabled *bool             `json:"security_policy_sync_enabled,omitempty"`
	ClusterType               string            `json:"cluster_type,omitempty"`
	ExternalNetworkId         string            `json:"external_network_id,omitempty"`
}

type UpdateOp string

const (
	AddOp     UpdateOp = "add"
	RemoveOp  UpdateOp = "remove"
	ReplaceOp UpdateOp = "replace"
)

type UpdateOpts struct {
	Op    UpdateOp    `json:"op" required:"true"`
	Path  string      `json:"path" required:"true"`
	Value interface{} `json:"value,omitempty"`
}

type ActionsBaseOpts struct {
	Action  string      `json:"action" required:"true"`
	Payload interface{} `json:"payload,omitempty"`
}

type UpgradeOpts struct {
	ClusterTemplateID string `json:"cluster_template_id" required:"true"`
	RollingEnabled    bool   `json:"rolling_enabled"`
}

type Cluster struct {
	APIAddress                string             `json:"api_address"`
	ClusterTemplateID         string             `json:"cluster_template_id"`
	CreatedAt                 time.Time          `json:"created_at"`
	DiscoveryURL              string             `json:"discovery_url"`
	KeyPair                   string             `json:"keypair"`
	Labels                    map[string]string  `json:"labels"`
	Links                     []gophercloud.Link `json:"links"`
	MasterFlavorID            string             `json:"master_flavor_id"`
	MasterAddresses           []string           `json:"master_addresses"`
	MasterCount               int                `json:"master_count"`
	Name                      string             `json:"name"`
	ProjectID                 string             `json:"project_id"`
	StackID                   string             `json:"stack_id"`
	Status                    string             `json:"status"`
	NewStatus                 string             `json:"new_status"`
	StatusReason              string             `json:"status_reason"`
	UUID                      string             `json:"uuid"`
	UpdatedAt                 time.Time          `json:"updated_at"`
	UserID                    string             `json:"user_id"`
	NetworkID                 string             `json:"network_id"`
	SubnetID                  string             `json:"subnet_id"`
	PodsNetworkCidr           string             `json:"pods_network_cidr"`
	FloatingIPEnabled         bool               `json:"floating_ip_enabled"`
	APILBVIP                  string             `json:"api_lb_vip"`
	APILBFIP                  string             `json:"api_lb_fip"`
	IngressFloatingIP         string             `json:"ingress_floating_ip"`
	RegistryAuthPassword      string             `json:"registry_auth_password"`
	AvailabilityZone          string             `json:"availability_zone,omitempty"`
	AvailabilityZones         []string           `json:"availability_zones,omitempty"`
	LoadbalancerSubnetID      string             `json:"loadbalancer_subnet_id"`
	InsecureRegistries        []string           `json:"insecure_registries,omitempty"`
	DNSDomain                 string             `json:"dns_domain,omitempty"`
	SecurityPolicySyncEnabled *bool              `json:"security_policy_sync_enabled,omitempty"`
	ClusterType               string             `json:"cluster_type,omitempty"`
}

// Map builds request params.
func (opts *CreateOpts) Map() (map[string]interface{}, error) {
	cluster, err := gophercloud.BuildRequestBody(*opts, "")
	return cluster, err
}

// Map builds request params.
func (opts *UpdateOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	return body, err
}

// Map builds request params.
func (opts *ActionsBaseOpts) Map() (map[string]interface{}, error) {
	cluster, err := gophercloud.BuildRequestBody(*opts, "")
	return cluster, err
}

// Map builds request params.
func (opts *UpgradeOpts) Map() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(*opts, "")
	return body, err
}

type clusterConfigResult struct {
	commonResult
}

type clusterDeleteResult struct {
	gophercloud.ErrResult
}

// Extract parses result into params for cluster.
func (r clusterConfigResult) Extract() (*Cluster, error) {
	var s *Cluster
	err := r.ExtractInto(&s)
	return s, err
}

func Create(client *gophercloud.ServiceClient, opts OptsBuilder) (r clusters.CreateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(clustersURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

// Get gets cluster data from vkcs.
func Get(client *gophercloud.ServiceClient, id string) (r clusterConfigResult) {
	log.Printf("GET cluster %s", id)
	resp, err := client.Get(clusterURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func KubeConfigGet(client *gophercloud.ServiceClient, id string) (string, error) {
	resp, err := client.Get(kubeConfigURL(client, id), nil, &gophercloud.RequestOpts{
		OkCodes:          []int{200},
		KeepResponseBody: true,
	})
	if err != nil {
		return "", errutil.ErrorWithRequestID(err, resp.Header.Get(errutil.RequestIDHeader))
	}
	buf := bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return "", errutil.ErrorWithRequestID(err, resp.Header.Get(errutil.RequestIDHeader))
	}
	return buf.String(), nil
}

func Update(client *gophercloud.ServiceClient, id string, opts []OptsBuilder) (r UpdateResult) {
	var o []map[string]interface{}
	for _, opt := range opts {
		b, err := opt.Map()
		if err != nil {
			r.Err = err
			return r
		}
		o = append(o, b)
	}
	resp, err := client.Patch(clusterURL(client, id), o, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Upgrade(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r clusters.UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(upgradeURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func UpdateMasters(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r clusters.UpdateResult) {
	b, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionsURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func SwitchState(client *gophercloud.ServiceClient, id string, opts OptsBuilder) (r clusters.UpdateResult) {
	reqBody, err := opts.Map()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(actionsURL(client, id), reqBody, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r clusterDeleteResult) {
	resp, err := client.Delete(clusterURL(client, id), &gophercloud.RequestOpts{})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	r.Err = errutil.ErrorWithRequestID(r.Err, r.Header.Get(errutil.RequestIDHeader))
	return
}
