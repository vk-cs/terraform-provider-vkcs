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
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key used to decrypt data",
			},
			"ciphertext": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Encrypted data",
			},
			"plaintext": {
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

	key, ok := d.Get("key").(string)
	if !ok {
		return diag.Errorf("Error retrieving key from resource: %s", err)
	}

	ciphertext, ok := d.Get("ciphertext").(string)
	if !ok {
		return diag.Errorf("Error retrieving ciphertext from resource: %s", err)
	}

	plaintext, err := decrypt(kmsV1Client, key, ciphertext)
	if err != nil {
		return diag.Errorf("Error decrypting: %s", err)
	}

	hash := sha256.Sum256([]byte(plaintext))
	d.SetId(hex.EncodeToString(hash[:]))

	err = d.Set("plaintext", plaintext)
	if err != nil {
		return diag.Errorf("Error setting plaintext: %s", err)
	}

	return nil
}
