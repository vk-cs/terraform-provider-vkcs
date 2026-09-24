package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceSecretsList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretsListReadContext,
		Schema:      map[string]*schema.Schema{},
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

	d.Set("secrets", secrets)
	return nil
}
