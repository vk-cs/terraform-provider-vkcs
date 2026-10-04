package kms

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/keys"
)

func listKeys(client *gophercloud.ServiceClient) ([]string, error) {
	k, err := keys.List(client).Extract()
	if err != nil {
		return nil, err
	}
	return k, nil
}

func getKey(client *gophercloud.ServiceClient, name string) (keys.Key, error) {
	k, err := keys.Get(client, name).Extract()
	if err != nil {
		return keys.Key{}, err
	}
	return k, nil
}

func encrypt(client *gophercloud.ServiceClient, key string, plaintext string) (string, error) {
	r, err := keys.Encrypt(client, key, keys.EncryptOpts{
		Plaintext: plaintext,
	}).Extract()
	if err != nil {
		return "", err
	}
	return r, nil
}

func decrypt(client *gophercloud.ServiceClient, key string, ciphertext string) (string, error) {
	r, err := keys.Decrypt(client, key, keys.DecryptOpts{
		Ciphertext: ciphertext,
	}).Extract()
	if err != nil {
		return "", err
	}
	return r, nil
}
