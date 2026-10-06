package kms

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/secrets"
)

func listSecrets(client *gophercloud.ServiceClient) ([]string, error) {
	s, err := secrets.List(client).Extract()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func getSecret(client *gophercloud.ServiceClient, path string) (secrets.GetSecret, error) {
	s, err := secrets.Get(client, path).Extract()
	if err != nil {
		return secrets.GetSecret{}, err
	}
	return s, nil
}

func createOrUpdateSecret(client *gophercloud.ServiceClient, path string, opts secrets.CreateOrUpdateOpts) (secrets.CreateOrUpdateSecret, error) {
	s, err := secrets.CreateOrUpdate(client, path, opts).Extract()
	if err != nil {
		return secrets.CreateOrUpdateSecret{}, err
	}
	return s, nil
}

func deleteSecret(client *gophercloud.ServiceClient, path string) error {
	return secrets.Delete(client, path).Extract()
}
