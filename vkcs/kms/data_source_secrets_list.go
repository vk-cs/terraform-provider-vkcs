package kms

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceSecretsList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretsListReadContext,
		Schema: map[string]*schema.Schema{
			"secrets_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"secrets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSecretsListReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	secrets, err := listSecrets(kmsV1Client)
	if err != nil {
		return diag.Errorf("Error listing secrets: %s", err)
	}

	var builder strings.Builder
	for _, secret := range secrets {
		_, err := builder.WriteString(secret)
		if err != nil {
			return diag.Errorf("Error building id: %s", err)
		}
	}
	hash := sha256.Sum256([]byte(fmt.Appendf(nil, "%s:%d", builder.String(), len(secrets))))

	d.SetId(hex.EncodeToString(hash[:]))
	err = d.Set("secrets_count", len(secrets))
	if err != nil {
		return diag.Errorf("Error setting secrets_count: %s", err)
	}
	err = d.Set("secrets", secrets)
	if err != nil {
		return diag.Errorf("Error setting secrets: %s", err)
	}
	return nil
}
