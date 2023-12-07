package clients

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/terraform/auth"
	"github.com/gophercloud/utils/terraform/mutexkv"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	sdkdiag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/version"
)

const (
	maxRetriesCount         = 3
	requestsMaxRetriesCount = 3
	requestsRetryDelay      = 1 * time.Second
)

// Config is interface to work with configer calls
type Config interface {
	LoadAndValidate() error
	GetRegion() string
	GetTenantID() string
	ComputeV2Client(region string) (*gophercloud.ServiceClient, error)
	ImageV2Client(region string) (*gophercloud.ServiceClient, error)
	NetworkingV2Client(region string, sdn string) (*gophercloud.ServiceClient, error)
	PublicDNSV2Client(region string) (*gophercloud.ServiceClient, error)
	BlockStorageV3Client(region string) (*gophercloud.ServiceClient, error)
	KeyManagerV1Client(region string) (*gophercloud.ServiceClient, error)
	ContainerInfraV1Client(region string) (*gophercloud.ServiceClient, error)
	ContainerInfraAddonsV1Client(region string) (*gophercloud.ServiceClient, error)
	IdentityV3Client(region string) (*gophercloud.ServiceClient, error)
	DatabaseV1Client(region string) (*gophercloud.ServiceClient, error)
	SharedfilesystemV2Client(region string) (*gophercloud.ServiceClient, error)
	LoadBalancerV2Client(region string) (*gophercloud.ServiceClient, error)
	BackupV1Client(region string, tenantID string) (*gophercloud.ServiceClient, error)
	MLPlatformV1Client(region string) (*gophercloud.ServiceClient, error)
	GetMutex() *mutexkv.MutexKV
}

// configer uses openstackbase.Config as the base/foundation of this provider's
type configer struct {
	auth.Config
	ContainerInfraV1MicroVersion string
}

func getConfigParam(d *schema.ResourceData, key string, envKey string, defaultVal string) (param string) {
	tfAttr := d.Get(key)
	if tfAttr != nil {
		param = tfAttr.(string)
	}
	if param == "" {
		param = os.Getenv(envKey)
	}
	if param == "" {
		param = defaultVal
	}
	return param
}

func ConfigureSdkProvider(d *schema.ResourceData, terraformVersion string) (Config, sdkdiag.Diagnostics) {
	containerInfraV1MicroVersion := d.Get("cloud_containers_api_version").(string)
	if containerInfraV1MicroVersion == "" {
		containerInfraV1MicroVersion = CloudContainersAPIVersion
	}

	config := &configer{
		auth.Config{
			Username:         getConfigParam(d, "username", "OS_USERNAME", ""),
			Password:         getConfigParam(d, "password", "OS_PASSWORD", ""),
			TenantID:         getConfigParam(d, "project_id", "OS_PROJECT_ID", ""),
			Region:           getConfigParam(d, "region", "OS_REGION_NAME", DefaultRegionName),
			IdentityEndpoint: getConfigParam(d, "auth_url", "OS_AUTH_URL", DefaultIdentityEndpoint),
			UserDomainID:     getConfigParam(d, "user_domain_id", "OS_USER_DOMAIN_ID", ""),
			UserDomainName:   getConfigParam(d, "user_domain_name", "OS_USER_DOMAIN_NAME", DefaultUserDomainName),
			EndpointType:     os.Getenv("OS_INTERFACE"),
			AllowReauth:      true,
			MaxRetries:       maxRetriesCount,
			TerraformVersion: terraformVersion,
			SDKVersion:       meta.SDKVersionString(),
			MutexKV:          mutexkv.NewMutexKV(),
		},
		containerInfraV1MicroVersion,
	}

	if config.UserDomainID != "" {
		config.UserDomainName = ""
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, sdkdiag.FromErr(err)
	}

	config.OsClient.UserAgent.Prepend(fmt.Sprintf("VKCS Terraform Provider %s", version.ProviderVersion))
	config.OsClient.RetryFunc = retryFunc

	return config, nil
}

var _ Config = &configer{}

// GetRegion is implementation of GetRegion method
func (c *configer) GetRegion() string {
	return c.Region
}

func (c *configer) GetTenantID() string {
	return c.TenantID
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

// ContainerInfraV1Client is implementation of ContainerInfraV1Client method
func (c *configer) ContainerInfraAddonsV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(newContainerInfraAddonsV1, region, "magnum-addons")
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

func (c *configer) BackupV1Client(region string, tenantID string) (*gophercloud.ServiceClient, error) {
	client, err := c.CommonServiceClientInit(newBackupV1, region, "data-protect")
	client.Endpoint = fmt.Sprintf("%s%s/", client.Endpoint, tenantID)
	return client, err
}

func (c *configer) MLPlatformV1Client(region string) (*gophercloud.ServiceClient, error) {
	client, err := c.CommonServiceClientInit(newMLPlatformV1, region, "mlplatform")
	return client, err
}

func (c *configer) GetMutex() *mutexkv.MutexKV {
	return c.Config.MutexKV
}

func (c *configer) setDefaults() {
	if c.TerraformVersion == "" {
		// Terraform 0.12 introduced this field to the protocol
		// We can therefore assume that if it's missing it's 0.10 or 0.11
		c.TerraformVersion = "0.11+compatible"
	}
	if c.ContainerInfraV1MicroVersion == "" {
		c.ContainerInfraV1MicroVersion = CloudContainersAPIVersion
	}
	if c.Region == "" {
		c.Region = DefaultRegionName
	}
	if c.IdentityEndpoint == "" {
		c.IdentityEndpoint = DefaultIdentityEndpoint
	}
	if c.UserDomainName == "" {
		c.UserDomainName = DefaultUserDomainName
	}
	if c.UserDomainID != "" {
		c.UserDomainName = ""
	}

	c.AllowReauth = true
	c.MaxRetries = maxRetriesCount
	c.MutexKV = mutexkv.NewMutexKV()
}

func (c *configer) updateWithEnv() {
	if c.Username == "" {
		c.Username = os.Getenv("OS_USERNAME")
	}
	if c.Password == "" {
		c.Password = os.Getenv("OS_PASSWORD")
	}
	if c.TenantID == "" {
		c.TenantID = os.Getenv("OS_PROJECT_ID")
	}
	if c.Region == "" {
		c.Region = os.Getenv("OS_REGION_NAME")
	}
	if c.IdentityEndpoint == "" {
		c.IdentityEndpoint = os.Getenv("OS_AUTH_URL")
	}
	if c.UserDomainID == "" {
		c.UserDomainID = os.Getenv("OS_USER_DOMAIN_ID")
	}
	if c.UserDomainName == "" {
		c.UserDomainName = os.Getenv("OS_USER_DOMAIN_NAME")
	}
	if c.EndpointType == "" {
		c.EndpointType = os.Getenv("OS_INTERFACE")
	}
}

func ConfigureProvider(ctx context.Context, req provider.ConfigureRequest) (Config, diag.Diagnostics) {
	var diags diag.Diagnostics
	config := configer{}

	req.Config.GetAttribute(ctx, path.Root("auth_url"), &config.IdentityEndpoint)
	req.Config.GetAttribute(ctx, path.Root("username"), &config.Username)
	req.Config.GetAttribute(ctx, path.Root("password"), &config.Password)
	req.Config.GetAttribute(ctx, path.Root("project_id"), &config.TenantID)
	req.Config.GetAttribute(ctx, path.Root("user_domain_id"), &config.UserDomainID)
	req.Config.GetAttribute(ctx, path.Root("user_domain_name"), &config.UserDomainName)
	req.Config.GetAttribute(ctx, path.Root("region"), &config.Region)
	req.Config.GetAttribute(ctx, path.Root("cloud_containers_api_version"), &config.ContainerInfraV1MicroVersion)
	config.updateWithEnv()
	config.TerraformVersion = req.TerraformVersion

	config.setDefaults()

	if err := config.LoadAndValidate(); err != nil {
		diags.AddError("Config validation error", err.Error())
		return nil, diags
	}

	config.OsClient.UserAgent.Prepend(fmt.Sprintf("VKCS Terraform Provider %s", version.ProviderVersion))
	config.OsClient.RetryFunc = retryFunc

	return &config, diags
}

func ConfigureFromEnv(ctx context.Context) (Config, error) {
	config := &configer{}
	config.updateWithEnv()
	config.setDefaults()

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return config, nil
}

func retryFunc(context context.Context, method, url string, options *gophercloud.RequestOpts, err error, failCount uint) error {
	if failCount >= requestsMaxRetriesCount {
		return err
	}
	if errutil.Any(err, []int{500, 501, 502, 503, 504}) {
		time.Sleep(requestsRetryDelay)
		return nil
	}

	return err
}
