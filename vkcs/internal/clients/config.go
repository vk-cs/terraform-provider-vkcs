package clients

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/terraform/auth"
	"github.com/gophercloud/utils/terraform/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
)

const (
	maxRetriesCount         = 3
	requestsMaxRetriesCount = 3
	requestsRetryDelay      = 30 * time.Millisecond
)

// Config is interface to work with configer calls
type Config interface {
	LoadAndValidate() error
	GetRegion() string
	ComputeV2Client(region string) (*gophercloud.ServiceClient, error)
	ImageV2Client(region string) (*gophercloud.ServiceClient, error)
	NetworkingV2Client(region string, sdn string) (*gophercloud.ServiceClient, error)
	PublicDNSV2Client(region string) (*gophercloud.ServiceClient, error)
	BlockStorageV3Client(region string) (*gophercloud.ServiceClient, error)
	KeyManagerV1Client(region string) (*gophercloud.ServiceClient, error)
	ContainerInfraV1Client(region string) (*gophercloud.ServiceClient, error)
	IdentityV3Client(region string) (*gophercloud.ServiceClient, error)
	DatabaseV1Client(region string) (*gophercloud.ServiceClient, error)
	SharedfilesystemV2Client(region string) (*gophercloud.ServiceClient, error)
	LoadBalancerV2Client(region string) (*gophercloud.ServiceClient, error)
	GetMutex() *mutexkv.MutexKV
}

// configer uses openstackbase.Config as the base/foundation of this provider's
type configer struct {
	auth.Config
	ContainerInfraV1MicroVersion string
}

func ConfigureProvider(d *schema.ResourceData, terraformVersion string) (Config, diag.Diagnostics) {
	config := &configer{
		auth.Config{
			Username:         d.Get("username").(string),
			Password:         d.Get("password").(string),
			TenantID:         d.Get("project_id").(string),
			Region:           d.Get("region").(string),
			IdentityEndpoint: d.Get("auth_url").(string),
			UserDomainID:     d.Get("user_domain_id").(string),
			UserDomainName:   d.Get("user_domain_name").(string),
			AllowReauth:      true,
			MaxRetries:       maxRetriesCount,
			TerraformVersion: terraformVersion,
			SDKVersion:       meta.SDKVersionString(),
			MutexKV:          mutexkv.NewMutexKV(),
		},
		d.Get("cloud_containers_api_version").(string),
	}

	if config.UserDomainID != "" {
		config.UserDomainName = ""
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, diag.FromErr(err)
	}
	return config, nil
}

var _ Config = &configer{}

// GetRegion is implementation of GetRegion method
func (c *configer) GetRegion() string {
	return c.Region
}

func (c *configer) ComputeV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.ComputeV2Client(region)
}

func (c *configer) ImageV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.ImageV2Client(region)
}

func (c *configer) NetworkingV2Client(region string, sdn string) (*gophercloud.ServiceClient, error) {
	client, err := c.Config.NetworkingV2Client(region)
	if err != nil {
		return client, err
	}
	err = networking.SelectSDN(client, sdn)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *configer) PublicDNSV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(newPublicDNSV2, region, "publicdns")
}

func (c *configer) BlockStorageV3Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.BlockStorageV3Client(region)
}

func (c *configer) KeyManagerV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.KeyManagerV1Client(region)
}

// DatabaseV1Client is implementation of DatabaseV1Client method
func (c *configer) DatabaseV1Client(region string) (*gophercloud.ServiceClient, error) {
	client, clientErr := c.Config.DatabaseV1Client(region)
	client.ProviderClient.RetryFunc = func(context context.Context, method, url string, options *gophercloud.RequestOpts, err error, failCount uint) error {
		if failCount >= requestsMaxRetriesCount {
			return err
		}
		switch errType := err.(type) {
		case gophercloud.ErrDefault500, gophercloud.ErrDefault503:
			time.Sleep(requestsRetryDelay)
			return nil
		case gophercloud.ErrUnexpectedResponseCode:
			if errType.Actual == http.StatusGatewayTimeout {
				time.Sleep(requestsRetryDelay)
				return nil
			}
			return err
		default:
			return err
		}
	}
	return client, clientErr
}

// ContainerInfraV1Client is implementation of ContainerInfraV1Client method
func (c *configer) ContainerInfraV1Client(region string) (*gophercloud.ServiceClient, error) {
	client, err := c.Config.ContainerInfraV1Client(region)
	if err != nil {
		return client, err
	}
	client.MoreHeaders = map[string]string{
		"MCS-API-Version": fmt.Sprintf("container-infra %s", c.ContainerInfraV1MicroVersion),
	}
	return client, err
}

// IdentityV3Client is implementation of ContainerInfraV1Client method
func (c *configer) IdentityV3Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.IdentityV3Client(region)
}

func (c *configer) SharedfilesystemV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.SharedfilesystemV2Client(region)
}

func (c *configer) LoadBalancerV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.Config.LoadBalancerV2Client(region)
}

func (c *configer) GetMutex() *mutexkv.MutexKV {
	return c.Config.MutexKV
}
