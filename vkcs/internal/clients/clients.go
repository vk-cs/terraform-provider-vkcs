package clients

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud"
)

const (
	ContainerInfraAPIVersion = "1.28"
)

type ErrServiceNotFound struct {
	ServiceName string
}

func NewErrEndpointNotFound(serviceName string) error {
	return ErrServiceNotFound{ServiceName: serviceName}
}

func (e ErrServiceNotFound) Error() string {
	return fmt.Sprintf("%s is not supported in your region. Please contact support.", e.ServiceName)
}

func initClientOpts(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts, clientType string) (*gophercloud.ServiceClient, error) {
	sc := new(gophercloud.ServiceClient)
	eo.ApplyDefaults(clientType)
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return sc, err
	}
	sc.ProviderClient = client
	sc.Endpoint = url
	sc.Type = clientType
	return sc, nil
}

func newPublicDNSV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "publicdns")
	sc.ResourceBase = sc.Endpoint + "v2/"
	return sc, err
}

func newContainerInfraAddonsV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "manage-cluster-addons")
	sc.ResourceBase = sc.Endpoint + "v1/"
	return sc, err
}

func newBackupV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "data-protect")
	return sc, err
}

func newMLPlatformV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "mlplatform")
	sc.ResourceBase = sc.Endpoint + "v1_0/"
	return sc, err
}

func newICSV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	baseURL, err := findBaseMonitoringURL(client, eo, "ics")
	if err != nil {
		return nil, err
	}

	sc := new(gophercloud.ServiceClient)
	eo.ApplyDefaults("ics")
	sc.ProviderClient = client
	sc.Type = "ics"

	sc.Endpoint = fmt.Sprint(baseURL, "infra/ics/")
	sc.ResourceBase = sc.Endpoint + "v1/"

	return sc, nil
}

func newTemplaterV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	baseURL, err := findBaseMonitoringURL(client, eo, "templater")
	if err != nil {
		return nil, err
	}

	sc := new(gophercloud.ServiceClient)
	eo.ApplyDefaults("templater")
	sc.ProviderClient = client
	sc.Type = "templater"

	sc.Endpoint = fmt.Sprint(baseURL, "infra/templater/")
	sc.ResourceBase = sc.Endpoint + "v2/"

	return sc, nil
}

func findBaseMonitoringURL(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts, service string) (string, error) {
	eo.ApplyDefaults("data-protect")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		var errNotFound *gophercloud.ErrEndpointNotFound
		if errors.As(err, &errNotFound) {
			return "", NewErrEndpointNotFound("data-protect")
		}

		return "", err
	}

	baseURL, _, found := strings.Cut(url, "infra/")
	if !found {
		return "", NewErrEndpointNotFound(service)
	}

	return baseURL, nil
}

func newCDNV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "cdn")
	return sc, err
}

func newDataPlatform(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "dataplatform")
	return sc, err
}
