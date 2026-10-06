package kms

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/keys"
)

const (
	KeyFieldName            = "name"
	KeyFieldType            = "type"
	KeyFieldDeletionAllowed = "deletion_allowed"

	KeysListFieldKeysCount = "keys_count"
	KeysListFieldKeys      = "keys"

	KeyCryptFieldKey        = "key"
	KeyCryptFieldPlaintext  = "plaintext"
	KeyCryptFieldCiphertext = "ciphertext"
)

func listKeys(client *gophercloud.ServiceClient) ([]string, error) {
	k, err := keys.List(client).Extract()
	if err != nil {
		return nil, err
	}
	return k, nil
}

func getKey(client *gophercloud.ServiceClient, name string) (keys.GetKey, error) {
	k, err := keys.Get(client, name).Extract()
	if err != nil {
		return keys.GetKey{}, err
	}
	return k, nil
}

func createKey(client *gophercloud.ServiceClient, name string, opts keys.CreateOpts) (keys.CreateKey, error) {
	k, err := keys.Create(client, name, opts).Extract()
	if err != nil {
		return keys.CreateKey{}, err
	}
	return k, nil
}

func updateKey(client *gophercloud.ServiceClient, name string, opts keys.UpdateOpts) (keys.UpdateKey, error) {
	k, err := keys.Update(client, name, opts).Extract()
	if err != nil {
		return keys.UpdateKey{}, err
	}
	return k, nil
}

func deleteKey(client *gophercloud.ServiceClient, name string) error {
	return keys.Delete(client, name).Extract()
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
