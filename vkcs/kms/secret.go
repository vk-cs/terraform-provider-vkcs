package kms

import (
	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/secrets"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
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

// secretStateRefreshFunc polls existence and readiness, including unfinished operations.
func secretStateRefreshFunc(client *gophercloud.ServiceClient, path string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		secret, err := getSecret(client, path)
		if isOperationUnfinished(err) {
			return secret, kmsStatePending, nil
		}
		if errutil.IsNotFound(err) {
			return secret, kmsStateMissing, nil
		}
		if err != nil {
			return nil, "", err
		}
		return secret, kmsStateReady, nil
	}
}

func secretDeleteProtectionStateRefreshFunc(client *gophercloud.ServiceClient, path string, expected *bool) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		protection, err := getSecretDeleteProtection(client, path)
		if isOperationUnfinished(err) {
			return protection, kmsStatePending, nil
		}
		if errutil.IsNotFound(err) {
			return protection, kmsStateMissing, nil
		}
		if err != nil {
			return nil, "", err
		}
		if expected != nil && protection != *expected {
			return protection, kmsStatePending, nil
		}
		return protection, kmsStateReady, nil
	}
}
