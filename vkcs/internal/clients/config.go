package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/gophercloud/utils/terraform/auth"
	"github.com/gophercloud/utils/terraform/mutexkv"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/monitoring/templater"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

const (
	maxRetriesCount         = 3
	requestsMaxRetriesCount = 3
	requestsRetryDelay      = 1 * time.Second
)

var (
	_ Config = (*config)(nil)
)

type Config interface {
	GetRegion() string
	GetProjectID() string
	GetToken() string
	GetMutex() *mutexkv.MutexKV

	BackupV1Client(region string, tenantID string) (*gophercloud.ServiceClient, error)
	BlockStorageV3Client(region string) (*gophercloud.ServiceClient, error)
	CDNV1Client(region string) (*gophercloud.ServiceClient, error)
	ComputeV2Client(region string) (*gophercloud.ServiceClient, error)
	ContainerInfraAddonsV1Client(region string) (*gophercloud.ServiceClient, error)
	ContainerInfraV1Client(region string) (*gophercloud.ServiceClient, error)
	DatabaseV1Client(region string) (*gophercloud.ServiceClient, error)
	DataPlatformClient(region string) (*gophercloud.ServiceClient, error)
	IAMServiceUsersV1Client(region string) (*gophercloud.ServiceClient, error)
	ICSV1Client(region string) (*gophercloud.ServiceClient, error)
	IdentityV3Client(region string) (*gophercloud.ServiceClient, error)
	ImageV2Client(region string) (*gophercloud.ServiceClient, error)
	KeyManagerV1Client(region string) (*gophercloud.ServiceClient, error)
	LoadBalancerV2Client(region string) (*gophercloud.ServiceClient, error)
	MLPlatformV1Client(region string) (*gophercloud.ServiceClient, error)
	NetworkingV2Client(region string, sdn string) (*gophercloud.ServiceClient, error)
	PublicDNSV2Client(region string) (*gophercloud.ServiceClient, error)
	SharedFilesystemV2Client(region string) (*gophercloud.ServiceClient, error)
	TemplaterV2Client(region string, projectID string) (*gophercloud.ServiceClient, error)
}

type config struct {
	auth.Config

	envPrefix                    string
	containerInfraV1MicroVersion string
	skipAuth                     bool
}

func (c *config) GetRegion() string {
	return c.Region
}

func (c *config) GetProjectID() string {
	return c.TenantID
}

func (c *config) GetToken() string {
	return c.OsClient.TokenID
}

func (c *config) GetMutex() *mutexkv.MutexKV {
	return c.MutexKV
}

func (c *config) BackupV1Client(region string, projectID string) (*gophercloud.ServiceClient, error) {
	client, err := c.initClient(newBackupV1, region, "backup")
	client.Endpoint = fmt.Sprintf("%s%s/", client.Endpoint, projectID)
	return client, err
}

func (c *config) BlockStorageV3Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newBlockStorageV3, region, "block-storage")
}

func (c *config) CDNV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newCDNV1, region, "cdn")
}

func (c *config) ComputeV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newComputeV2, region, "compute")
}

func (c *config) ContainerInfraV1Client(region string) (*gophercloud.ServiceClient, error) {
	client, err := c.initClient(newContainerInfraV1, region, "container-infra")
	if err != nil {
		return client, err
	}

	client.MoreHeaders = map[string]string{
		"MCS-API-Version": fmt.Sprintf("container-infra %s", c.containerInfraV1MicroVersion),
	}

	return client, err
}

func (c *config) ContainerInfraAddonsV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newContainerInfraAddonsV1, region, "container-infra-addons")
}

func (c *config) DatabaseV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newDatabaseV1, region, "database")
}

func (c *config) DataPlatformClient(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newDataPlatform, region, "data-platform")
}

func (c *config) IAMServiceUsersV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newIAMServiceUsersV1, region, "iam-service-users")
}

func (c *config) ICSV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newICSV1, region, "ics")
}

func (c *config) IdentityV3Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newIdentityV3, region, "identity")
}

func (c *config) ImageV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newImageV2, region, "image")
}

func (c *config) KeyManagerV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newKeyManagerV1, region, "key-manager")
}

func (c *config) LoadBalancerV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newLoadBalancerV2, region, "load-balancer")
}

func (c *config) MLPlatformV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newMLPlatformV1, region, "mlplatform")
}

func (c *config) NetworkingV2Client(region string, sdn string) (*gophercloud.ServiceClient, error) {
	client, err := c.initClient(newNetworkV2, region, "networking")
	if err != nil {
		return client, err
	}

	err = networking.SelectSDN(client, sdn)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *config) PublicDNSV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newPublicDNSV2, region, "public-dns")
}

func (c *config) SharedFilesystemV2Client(region string) (*gophercloud.ServiceClient, error) {
	return c.initClient(newSharedFilesystemV2, region, "shared-filesystem")
}

func (c *config) TemplaterV2Client(region string, projectID string) (*gophercloud.ServiceClient, error) {
	client, err := c.initClient(newTemplaterV2, region, "templater")
	if err != nil {
		return nil, err
	}

	if err = templater.ListUsers(client, projectID).ExtractErr(); err != nil {
		if errutil.IsNotFound(err) {
			return nil, newErrEndpointNotFound("templater")
		}

		return nil, err
	}

	return client, nil
}

type clientFactoryFn func(*gophercloud.ProviderClient, clientOpts) (*gophercloud.ServiceClient, error)

func (c *config) initClient(newClient clientFactoryFn, region, service string) (*gophercloud.ServiceClient, error) {
	endpointOverride := c.determineEndpoint(service)

	if !c.skipAuth {
		if err := c.Authenticate(); err != nil {
			return nil, err
		}
	} else {
		if endpointOverride == "" {
			return nil, fmt.Errorf("endpoint override for `%s` is not provided", service)
		}

		c.OsClient.SetToken(c.Token)
	}

	opts := clientOpts{
		EndpointOpts: gophercloud.EndpointOpts{
			Region:       c.DetermineRegion(region),
			Availability: clientconfig.GetEndpointType(c.EndpointType),
		},
		EndpointOverride: endpointOverride,
	}

	client, err := newClient(c.OsClient, opts)
	if err != nil {
		return client, err
	}

	return client, nil
}

func (c *config) determineEndpoint(service string) string {
	if service == "identity" {
		return c.IdentityEndpoint
	}

	if v, ok := c.EndpointOverrides[service]; ok {
		return v.(string)
	}

	return getEnv(c.envPrefix, fmt.Sprintf("%s_ENDPOINT_OVERRIDE", service))
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
