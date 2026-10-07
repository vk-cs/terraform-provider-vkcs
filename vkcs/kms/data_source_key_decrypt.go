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

func DataSourceKeyDecrypt() *schema.Resource {
	return &schema.Resource{
		Description: "Data source containing decrypt result by KMS key",
		ReadContext: dataSourceKeyDecryptReadContext,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			KeyCryptFieldKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key used to decrypt data",
			},
			KeyCryptFieldCiphertext: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Encrypted data",
			},
			KeyCryptFieldPlaintext: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Resulting decrypted data encoded in base64",
			},
		},
	}
}

func dataSourceKeyDecryptReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	key, ok := d.Get(KeyCryptFieldKey).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, KeyCryptFieldKey)
	}

	ciphertext, ok := d.Get(KeyCryptFieldCiphertext).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, KeyCryptFieldCiphertext)
	}

	plaintext, err := decrypt(kmsV1Client, key, ciphertext)
	if err != nil {
		return diag.Errorf("Error decrypting: %s", err)
	}

	hash := sha256.Sum256([]byte(plaintext))
	d.SetId(hex.EncodeToString(hash[:]))

	err = d.Set(KeyCryptFieldPlaintext, plaintext)
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, KeyCryptFieldPlaintext, err)
	}

	return nil
}
