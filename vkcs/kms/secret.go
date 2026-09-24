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

func getSecret(client *gophercloud.ServiceClient, name string) (secrets.Secret, error) {
	s, err := secrets.Get(client, name).Extract()
	if err != nil {
		return secrets.Secret{}, err
	}
	return s, nil
}
