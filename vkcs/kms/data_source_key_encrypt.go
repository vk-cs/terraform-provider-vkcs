package kms

import (
	"context"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceKeyEncrypt() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeyEncryptReadContext,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plaintext": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ciphertext": {
				Type: schema.TypeString,
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

	key, ok := d.Get("key").(string)
	if !ok {
		return diag.Errorf("Error retrieving key from resource: %s", err)
	}

	plaintext, ok := d.Get("plaintext").(string)
	if !ok {
		return diag.Errorf("Error retrieving name from resource: %s", err)
	}

	ciphertext, err := encrypt(kmsV1Client, key, plaintext)
	if err != nil {
		return diag.Errorf("Error encrypting: %s", err)
	}

	d.SetId(base64.StdEncoding.EncodeToString([]byte(ciphertext)))
	err = d.Set("ciphertext", ciphertext)
	if err != nil {
		return diag.Errorf("Error setting ciphertext: %s", err)
	}

	return nil
}
