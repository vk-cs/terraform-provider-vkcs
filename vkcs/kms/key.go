package kms

import (
	"github.com/gophercloud/gophercloud"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/keys"
)

func encrypt(client *gophercloud.ServiceClient, key string, plaintext string) (string, error) {
	r, err := keys.Encrypt(client, key, keys.EncryptOpts{
		Plaintext: plaintext,
	}).Extract()
	if err != nil {
		return "", err
	}
	return r, nil
}
