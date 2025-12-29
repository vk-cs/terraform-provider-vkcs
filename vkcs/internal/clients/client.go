package clients

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud"
)

const (
	computeAPIMicroVersion = "2.42"
)

type errServiceNotFound struct {
	ServiceName string
}

func newErrEndpointNotFound(serviceName string) error {
	return errServiceNotFound{ServiceName: serviceName}
}

func (e errServiceNotFound) Error() string {
	return fmt.Sprintf("%s is not supported in your region. Please contact support.", e.ServiceName)
}

type clientOpts struct {
	gophercloud.EndpointOpts

	EndpointOverride string
}

func newBackupV1(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	return initClientOptsNew(client, opts, "data-protect")
}

func newBlockStorageV3(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	return initClientOptsNew(client, opts, "volumev3")
}

func newCDNV1(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	return initClientOptsNew(client, opts, "cdn")
}

func newComputeV2(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	c, err := initClientOptsNew(client, opts, "compute")
	if err != nil {
		return c, err
	}

	c.Microversion = computeAPIMicroVersion

	return c, nil
}

func newContainerInfraV1(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	return initClientOptsNew(client, opts, "container-infra")
}

func newContainerInfraAddonsV1(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOptsNew(client, opts, "manage-cluster-addons")
	sc.ResourceBase = sc.Endpoint + "v1/"
	return sc, err
}

func newDatabaseV1(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	return initClientOptsNew(client, opts, "database")
}

func newDataPlatform(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	return initClientOptsNew(client, opts, "dataplatform")
}

func newIAMServiceUsersV1(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	return initClientOptsNew(client, opts, "service-users")
}

func newICSV1(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	if opts.EndpointOverride == "" {
		baseURL, err := findBaseMonitoringURL(client, opts, "ics")
		if err != nil {
			return nil, err
		}
		opts.EndpointOverride = fmt.Sprint(baseURL, "infra/ics/")
	}

	sc, err := initClientOptsNew(client, opts, "ics")
	if err != nil {
		return nil, err
	}

	sc.ResourceBase = sc.Endpoint + "v1/"

	return sc, nil
}

func newIdentityV3(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	if opts.EndpointOverride == "" {
		return nil, newErrEndpointNotFound("identity")
	}

	return initClientOptsNew(client, opts, "identity")
}

func newImageV2(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOptsNew(client, opts, "image")
	sc.ResourceBase = sc.Endpoint + "v2/"
	return sc, err
}

func newKeyManagerV1(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOptsNew(client, opts, "key-manager")
	sc.ResourceBase = sc.Endpoint + "v1/"
	return sc, err
}

func newLoadBalancerV2(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOptsNew(client, opts, "load-balancer")
	endpoint := strings.ReplaceAll(sc.Endpoint, "v2.0/", "")
	sc.ResourceBase = endpoint + "v2.0/"
	return sc, err
}

func newMLPlatformV1(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOptsNew(client, opts, "mlplatform")
	sc.ResourceBase = sc.Endpoint + "v1_0/"
	return sc, err
}

func newNetworkV2(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOptsNew(client, opts, "network")
	sc.ResourceBase = sc.Endpoint + "v2.0/"
	return sc, err
}

func newPublicDNSV2(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOptsNew(client, opts, "publicdns")
	sc.ResourceBase = sc.Endpoint + "v2/"
	return sc, err
}

func newSharedFilesystemV2(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	return initClientOptsNew(client, opts, "sharev2")
}

func newTemplaterV2(client *gophercloud.ProviderClient, opts clientOpts) (*gophercloud.ServiceClient, error) {
	if opts.EndpointOverride == "" {
		baseURL, err := findBaseMonitoringURL(client, opts, "templater")
		if err != nil {
			return nil, err
		}
		opts.EndpointOverride = fmt.Sprint(baseURL, "infra/templater/")
	}

	sc, err := initClientOptsNew(client, opts, "templater")
	if err != nil {
		return nil, err
	}

	sc.ResourceBase = sc.Endpoint + "v2/"

	return sc, nil
}

func initClientOptsNew(client *gophercloud.ProviderClient, opts clientOpts, clientType string) (*gophercloud.ServiceClient, error) {
	sc := new(gophercloud.ServiceClient)
	sc.ProviderClient = client
	opts.ApplyDefaults(clientType)
	sc.Type = clientType

	if eo := opts.EndpointOverride; eo != "" {
		sc.Endpoint = eo
		sc.ResourceBase = ""
	} else if sc.Endpoint == "" {
		if client.EndpointLocator == nil {
			return nil, fmt.Errorf("endpoint locator for client `%s` is not provided", clientType)
		}

		var err error
		sc.Endpoint, err = client.EndpointLocator(opts.EndpointOpts)
		if err != nil {
			return sc, err
		}
	}

	return sc, nil
}

func findBaseMonitoringURL(client *gophercloud.ProviderClient, opts clientOpts, service string) (string, error) {
	opts.ApplyDefaults("data-protect")
	url, err := client.EndpointLocator(opts.EndpointOpts)
	if err != nil {
		var errNotFound *gophercloud.ErrEndpointNotFound
		if errors.As(err, &errNotFound) {
			return "", newErrEndpointNotFound("data-protect")
		}

		return "", err
	}

	baseURL, _, found := strings.Cut(url, "infra/")
	if !found {
		return "", newErrEndpointNotFound(service)
	}

	return baseURL, nil
}
