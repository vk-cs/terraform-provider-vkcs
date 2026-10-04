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

func DataSourceKeysList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeysListReadContext,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"keys_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"keys": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceKeysListReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	keys, err := listKeys(kmsV1Client)
	if err != nil {
		return diag.Errorf("Error listing secrets: %s", err)
	}

	var builder strings.Builder
	for _, secret := range keys[:min(len(keys), maxItemsInListForIDBuild)] {
		_, err := builder.WriteString(secret)
		if err != nil {
			return diag.Errorf("Error building id: %s", err)
		}
	}
	hash := sha256.Sum256(fmt.Appendf(nil, "%s:%d", builder.String(), len(keys)))

	d.SetId(hex.EncodeToString(hash[:]))
	err = d.Set("keys_count", len(keys))
	if err != nil {
		return diag.Errorf("Error setting keys_count: %s", err)
	}
	err = d.Set("keys", keys)
	if err != nil {
		return diag.Errorf("Error setting keys: %s", err)
	}
	return nil
}
