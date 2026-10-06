package kms

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceKeyEncrypt() *schema.Resource {
	return &schema.Resource{
		Description: "Data source containing encrypt result by KMS key",
		ReadContext: dataSourceKeyEncryptReadContext,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			KeyCryptFieldKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key used to encrypt data",
			},
			KeyCryptFieldPlaintext: {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Data to encrypt encoded in base64",
			},
			KeyCryptFieldCiphertext: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resulting encrypted data",
			},
		},
	}
}

func dataSourceKeyEncryptReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	key, ok := d.Get(KeyCryptFieldKey).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, KeyCryptFieldKey)
	}

	plaintext, ok := d.Get(KeyCryptFieldPlaintext).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, KeyCryptFieldPlaintext)
	}

	ciphertext, err := encrypt(kmsV1Client, key, plaintext)
	if err != nil {
		return diag.Errorf("Error encrypting: %s", err)
	}

	hash := sha256.Sum256([]byte(ciphertext))
	d.SetId(hex.EncodeToString(hash[:]))

	err = d.Set(KeyCryptFieldCiphertext, ciphertext)
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, KeyCryptFieldCiphertext, err)
	}

	return nil
}
