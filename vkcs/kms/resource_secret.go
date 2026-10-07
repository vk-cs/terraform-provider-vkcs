package kms

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/secrets"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func ResourceSecret() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource representing KMS secret",
		CreateContext: resourceSecretCreateContext,
		ReadContext:   resourceSecretReadContext,
		UpdateContext: resourceSecretUpdateContext,
		DeleteContext: resourceSecretDeleteContext,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			SecretFieldPath: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identifier (path) of the secret within the KMS secret store. Changing this forces a new resource.",
			},
			SecretFieldDataJSON: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Data stored in the secret as a JSON object, represented as a string. Use `jsonencode` to encode a Terraform object. Refreshed from KMS when the resource is read.",
			},
			SecretFieldCreatedTime: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Creation time returned in the secret metadata, formatted as a Go time string, for example `2026-09-01 10:00:00 +0000 UTC`.",
			},
			SecretFieldVersion: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Current version of the secret returned in its metadata.",
			},
			SecretFieldDeleteProtection: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the secret is protected from deletion. Defaults to `true`. Set to `false` and apply before destroying the secret.",
			},
		},
	}
}

func resourceSecretCreateContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	path, ok := d.Get(SecretFieldPath).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, SecretFieldPath)
	}

	data, ok := d.Get(SecretFieldDataJSON).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, SecretFieldDataJSON)
	}

	secretData := make(map[string]any)
	err = json.Unmarshal([]byte(data), &secretData)
	if err != nil {
		return diag.Errorf("Error unmarshalling data: %s", err)
	}

	secret, err := createOrUpdateSecret(kmsV1Client, path, secrets.CreateOrUpdateOpts{
		Data: secretData,
	})
	if err != nil {
		return diag.Errorf("Error creating secret: %s", err)
	}

	deleteProtection, ok := d.Get(SecretFieldDeleteProtection).(bool)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, SecretFieldDeleteProtection)
	}

	// We need to wait for the secret to be ready, then set this flag
	for range retriesCount {
		err = setSecretDeleteProtection(kmsV1Client, path, secrets.SetDeleteProtectionOpts{
			DeleteProtection: deleteProtection,
		})
		if err == nil {
			break
		}
		time.Sleep(defaultThreshold)
	}
	if err != nil {
		return diag.Errorf("Error updating secret delete protection: %s", err)
	}

	d.SetId(path)

	err = d.Set(SecretFieldCreatedTime, secret.Data.CreatedTime.String())
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, SecretFieldCreatedTime, err)
	}

	err = d.Set(SecretFieldVersion, secret.Data.Version)
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, SecretFieldVersion, err)
	}

	return nil
}

func resourceSecretReadContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	path, ok := d.Get(SecretFieldPath).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, SecretFieldPath)
	}

	secret, err := getSecret(kmsV1Client, path)
	if err != nil {
		return diag.Errorf("Error getting secret: %s", err)
	}

	secretString, err := json.Marshal(secret.Data.Data)
	if err != nil {
		return diag.Errorf("Error creating data_json: %s", err)
	}

	deleteProtection, err := getSecretDeleteProtection(kmsV1Client, path)
	if err != nil {
		return diag.Errorf("Error getting secret delete protection: %s", err)
	}

	err = d.Set(SecretFieldDataJSON, string(secretString))
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, SecretFieldDataJSON, err)
	}

	err = d.Set(SecretFieldCreatedTime, secret.Data.Metadata.CreatedTime.String())
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, SecretFieldCreatedTime, err)
	}

	err = d.Set(SecretFieldVersion, secret.Data.Metadata.Version)
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, SecretFieldVersion, err)
	}

	err = d.Set(SecretFieldDeleteProtection, deleteProtection)
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, SecretFieldDeleteProtection, err)
	}

	return nil
}

func resourceSecretUpdateContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	path, ok := d.Get(SecretFieldPath).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, SecretFieldPath)
	}

	data, ok := d.Get(SecretFieldDataJSON).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, SecretFieldDataJSON)
	}

	secretData := make(map[string]any)
	err = json.Unmarshal([]byte(data), &secretData)
	if err != nil {
		return diag.Errorf("Error unmarshalling data: %s", err)
	}

	secret, err := createOrUpdateSecret(kmsV1Client, path, secrets.CreateOrUpdateOpts{
		Data: secretData,
	})
	if err != nil {
		return diag.Errorf("Error creating secret: %s", err)
	}

	deleteProtection, ok := d.Get(SecretFieldDeleteProtection).(bool)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, SecretFieldDeleteProtection)
	}

	// We need to wait for the secret to be ready, then set this flag
	for range retriesCount {
		err = setSecretDeleteProtection(kmsV1Client, path, secrets.SetDeleteProtectionOpts{
			DeleteProtection: deleteProtection,
		})
		if err == nil {
			break
		}
		time.Sleep(defaultThreshold)
	}
	if err != nil {
		return diag.Errorf("Error updating secret delete protection: %s", err)
	}

	err = d.Set(SecretFieldCreatedTime, secret.Data.CreatedTime.String())
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, SecretFieldCreatedTime, err)
	}

	err = d.Set(SecretFieldVersion, secret.Data.Version)
	if err != nil {
		return diag.Errorf(diagSetErrorTemplate, SecretFieldVersion, err)
	}

	return nil
}

func resourceSecretDeleteContext(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(clients.Config)
	kmsV1Client, err := config.KMSV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KMS client: %s", err)
	}

	path, ok := d.Get(SecretFieldPath).(string)
	if !ok {
		return diag.Errorf(diagRetrieveErrorTemplate, SecretFieldPath)
	}

	err = deleteSecret(kmsV1Client, path)
	if err != nil {
		return diag.Errorf("Error deleting secret: %s", err)
	}

	// Waiting for the secret to be fully deleted
	for range retriesCount {
		var err error
		_, err = getSecret(kmsV1Client, path)
		if err != nil {
			break
		}
		time.Sleep(defaultThreshold)
	}

	return nil
}
