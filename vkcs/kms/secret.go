package kms

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/secrets"
)

const (
	SecretFieldPath             = "path"
	SecretFieldDataJSON         = "data_json"
	SecretFieldCreatedTime      = "created_time"
	SecretFieldVersion          = "version"
	SecretFieldDeleteProtection = "delete_protection"

	SecretsListFieldSecretsCount = "secrets_count"
	SecretsListFieldSecrets      = "secrets"
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

func getSecretDeleteProtection(client *gophercloud.ServiceClient, path string) (bool, error) {
	deleteProtection, err := secrets.GetDeleteProtection(client, path).Extract()
	if err != nil {
		return false, err
	}
	return deleteProtection, nil
}

func createOrUpdateSecret(client *gophercloud.ServiceClient, path string, opts secrets.CreateOrUpdateOpts) (secrets.CreateOrUpdateSecret, error) {
	s, err := secrets.CreateOrUpdate(client, path, opts).Extract()
	if err != nil {
		return secrets.CreateOrUpdateSecret{}, err
	}
	return s, nil
}

func setSecretDeleteProtection(client *gophercloud.ServiceClient, path string, opts secrets.SetDeleteProtectionOpts) error {
	return secrets.SetDeleteProtection(client, path, opts).Extract()
}

func deleteSecret(client *gophercloud.ServiceClient, path string) error {
	return secrets.Delete(client, path).Extract()
}
