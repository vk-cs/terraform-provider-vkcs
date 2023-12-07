package clients

import "github.com/gophercloud/gophercloud"

const (
	ContainerInfraAPIVersion = "1.28"
)

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
