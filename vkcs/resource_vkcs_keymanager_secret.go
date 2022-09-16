package vkcs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
)

func resourceKeyManagerSecret() *schema.Resource {
	ret := &schema.Resource{
		CreateContext: resourceKeyManagerSecretCreate,
		ReadContext:   resourceKeyManagerSecretRead,
		UpdateContext: resourceKeyManagerSecretUpdate,
		DeleteContext: resourceKeyManagerSecretDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the KeyManager client. A KeyManager client is needed to create a secret. If omitted, the `region` argument of the provider is used. Changing this creates a new V1 secret.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Human-readable name for the Secret. Does not have to be unique.",
			},

			"bit_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Metadata provided by a user or system for informational purposes.",
			},

			"algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Metadata provided by a user or system for informational purposes.",
			},

			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the secret.",
			},

			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Metadata provided by a user or system for informational purposes.",
			},

			"secret_ref": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret reference / where to find the secret.",
			},

			"secret_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"symmetric", "public", "private", "passphrase", "certificate", "opaque",
				}, false),
				Description: "Used to indicate the type of secret being stored. For more information see [Secret types](https://docs.openstack.org/barbican/latest/api/reference/secret_types.html).",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the secret.",
			},

			"payload": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
				Computed:  true,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					return strings.TrimSpace(o) == strings.TrimSpace(n)
				},
				Description: "The secret's data to be stored. **payload\\_content\\_type** must also be supplied if **payload** is included.",
			},

			"payload_content_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"text/plain", "text/plain;charset=utf-8", "text/plain; charset=utf-8", "application/octet-stream", "application/pkcs8",
				}, true),
				Description: "(required if **payload** is included) The media type for the content of the payload. Must be one of `text/plain`, `text/plain;charset=utf-8`, `text/plain; charset=utf-8`, `application/octet-stream`, `application/pkcs8`.",
			},

			"payload_content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"base64", "binary",
				}, false),
				Description: "(required if **payload** is encoded) The encoding used for the payload to be able to include it in the JSON request. Must be either `base64` or `binary`.",
			},

			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    false,
				Description: "Additional Metadata for the secret.",
			},

			"acl": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Allows to control an access to a secret. Currently only the `read` operation is supported. If not specified, the secret is accessible project wide.",
			},

			"expiration": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
				Description:  "The expiration time of the secret in the RFC3339 timestamp format (e.g. `2019-03-09T12:58:49Z`). If omitted, a secret will never expire. Changing this creates a new secret.",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the secret ACL was created.",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the secret ACL was last updated.",
			},

			"content_types": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The map of the content types, assigned on the secret.",
			},

			"all_metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The map of metadata, assigned on the secret, which has been explicitly and implicitly added.",
			},
		},
		Description: "Manages a key secret resource within VKCS.\n\n~> **Important Security Notice** The payload of this resource will be stored *unencrypted* in your Terraform state file. **Use of this resource for production deployments is *not* recommended**. [Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).",

		CustomizeDiff: customdiff.Sequence(
			// Clear the diff if the source payload is base64 encoded.
			func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return resourceSecretV1PayloadBase64CustomizeDiff(diff)
			},
		),
	}

	elem := &schema.Resource{
		Schema: make(map[string]*schema.Schema),
	}
	for _, aclOp := range getSupportedACLOperations() {
		elem.Schema[aclOp] = getACLSchema()
		elem.Schema[aclOp].Description = fmt.Sprintf("Block that describes %s operation.", aclOp)
	}
	ret.Schema["acl"].Elem = elem

	return ret
}

func resourceKeyManagerSecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	kmClient, err := config.KeyManagerV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KeyManager client: %s", err)
	}

	var expiration *time.Time
	if v, err := time.Parse(time.RFC3339, d.Get("expiration").(string)); err == nil {
		expiration = &v
	}

	secretType := keyManagerSecretSecretType(d.Get("secret_type").(string))

	createOpts := secrets.CreateOpts{
		Name:       d.Get("name").(string),
		Algorithm:  d.Get("algorithm").(string),
		BitLength:  d.Get("bit_length").(int),
		Mode:       d.Get("mode").(string),
		Expiration: expiration,
		SecretType: secretType,
	}

	log.Printf("[DEBUG] Create Options for resource_keymanager_secret_v1: %#v", createOpts)

	var secret *secrets.Secret
	secret, err = secrets.Create(kmClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_keymanager_secret: %s", err)
	}

	uuid := keyManagerSecretGetUUIDfromSecretRef(secret.SecretRef)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"ACTIVE"},
		Refresh:    keyManagerSecretWaitForSecretCreation(kmClient, uuid),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_keymanager_secret: %s", err)
	}

	d.SetId(uuid)

	d.Partial(true)

	// set the acl first before uploading the payload
	if acl, ok := d.GetOk("acl"); ok {
		setOpts := expandKeyManagerACLs(acl)
		_, err = acls.SetSecretACL(kmClient, uuid, setOpts).Extract()
		if err != nil {
			return diag.Errorf("Error settings ACLs for the vkcs_keymanager_secret: %s", err)
		}
	}

	// set the payload
	updateOpts := secrets.UpdateOpts{
		Payload:         d.Get("payload").(string),
		ContentType:     d.Get("payload_content_type").(string),
		ContentEncoding: d.Get("payload_content_encoding").(string),
	}
	err = secrets.Update(kmClient, uuid, updateOpts).Err
	if err != nil {
		return diag.Errorf("Error setting vkcs_keymanager_secret payload: %s", err)
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_keymanager_secret: %s", err)
	}

	// set the metadata
	var metadataCreateOpts secrets.MetadataOpts = flattenKeyManagerSecretMetadata(d)

	log.Printf("[DEBUG] Metadata Create Options for resource_keymanager_secret_metadata_v1 %s: %#v", uuid, metadataCreateOpts)

	if len(metadataCreateOpts) > 0 {
		_, err = secrets.CreateMetadata(kmClient, uuid, metadataCreateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error creating metadata for vkcs_keymanager_secret with ID %s: %s", uuid, err)
		}

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"PENDING"},
			Target:     []string{"ACTIVE"},
			Refresh:    keyManagerSecretMetadataV1WaitForSecretMetadataCreation(kmClient, uuid),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error creating metadata for vkcs_keymanager_secret %s: %s", uuid, err)
		}
	}

	d.Partial(false)

	return resourceKeyManagerSecretRead(ctx, d, meta)
}

func resourceKeyManagerSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	kmClient, err := config.KeyManagerV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS keymanager client: %s", err)
	}

	secret, err := secrets.Get(kmClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_keymanager_secret"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_keymanager_secret %s: %#v", d.Id(), secret)

	d.Set("name", secret.Name)

	d.Set("bit_length", secret.BitLength)
	d.Set("algorithm", secret.Algorithm)
	d.Set("creator_id", secret.CreatorID)
	d.Set("mode", secret.Mode)
	d.Set("secret_ref", secret.SecretRef)
	d.Set("secret_type", secret.SecretType)
	d.Set("status", secret.Status)
	d.Set("created_at", secret.Created.Format(time.RFC3339))
	d.Set("updated_at", secret.Updated.Format(time.RFC3339))
	d.Set("content_types", secret.ContentTypes)

	// don't fail, if the default key doesn't exist
	payloadContentType := secret.ContentTypes["default"]
	d.Set("payload_content_type", payloadContentType)

	d.Set("payload", keyManagerSecretGetPayload(kmClient, d.Id()))
	metadataMap, err := secrets.GetMetadata(kmClient, d.Id()).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get %s secret metadata: %s", d.Id(), err)
	}
	d.Set("all_metadata", metadataMap)

	if secret.Expiration == (time.Time{}) {
		d.Set("expiration", "")
	} else {
		d.Set("expiration", secret.Expiration.Format(time.RFC3339))
	}

	acl, err := acls.GetSecretACL(kmClient, d.Id()).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get %s secret acls: %s", d.Id(), err)
	}
	d.Set("acl", flattenKeyManagerACLs(acl))

	// Set the region
	d.Set("region", getRegion(d, config))

	return nil
}

func resourceKeyManagerSecretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	kmClient, err := config.KeyManagerV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS keymanager client: %s", err)
	}

	if d.HasChange("acl") {
		updateOpts := expandKeyManagerACLs(d.Get("acl"))
		_, err := acls.UpdateSecretACL(kmClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating vkcs_keymanager_secret %s acl: %s", d.Id(), err)
		}
	}

	if d.HasChange("metadata") {
		var metadataToDelete []string
		var metadataToAdd []string
		var metadataToUpdate []string

		o, n := d.GetChange("metadata")
		oldMetadata := o.(map[string]interface{})
		newMetadata := n.(map[string]interface{})
		existingMetadata := d.Get("all_metadata").(map[string]interface{})

		// Determine if any metadata keys were removed from the configuration.
		// Then request those keys to be deleted.
		for oldKey := range oldMetadata {
			if _, ok := newMetadata[oldKey]; !ok {
				metadataToDelete = append(metadataToDelete, oldKey)
			}
		}

		log.Printf("[DEBUG] Deleting the following items from metadata for vkcs_keymanager_secret %s: %v", d.Id(), metadataToDelete)

		for _, key := range metadataToDelete {
			err := secrets.DeleteMetadatum(kmClient, d.Id(), key).ExtractErr()
			if err != nil {
				return diag.Errorf("Error deleting vkcs_keymanager_secret %s metadata %s: %s", d.Id(), key, err)
			}
		}

		// Determine if any metadata keys were updated or added in the configuration.
		// Then request those keys to be updated or added.
		for newKey, newValue := range newMetadata {
			if oldValue, ok := oldMetadata[newKey]; ok {
				if newValue != oldValue {
					metadataToUpdate = append(metadataToUpdate, newKey)
				}
			} else if existingValue, ok := existingMetadata[newKey]; ok {
				if newValue != existingValue {
					metadataToUpdate = append(metadataToUpdate, newKey)
				}
			} else {
				metadataToAdd = append(metadataToAdd, newKey)
			}
		}

		log.Printf("[DEBUG] Updating the following items in metadata for vkcs_keymanager_secret %s: %v", d.Id(), metadataToUpdate)

		for _, key := range metadataToUpdate {
			var metadatumOpts secrets.MetadatumOpts
			metadatumOpts.Key = key
			metadatumOpts.Value = newMetadata[key].(string)
			_, err := secrets.UpdateMetadatum(kmClient, d.Id(), metadatumOpts).Extract()
			if err != nil {
				return diag.Errorf("Error updating vkcs_keymanager_secret %s metadata %s: %s", d.Id(), key, err)
			}
		}

		log.Printf("[DEBUG] Adding the following items to metadata for vkcs_keymanager_secret %s: %v", d.Id(), metadataToAdd)

		for _, key := range metadataToAdd {
			var metadatumOpts secrets.MetadatumOpts
			metadatumOpts.Key = key
			metadatumOpts.Value = newMetadata[key].(string)
			err := secrets.CreateMetadatum(kmClient, d.Id(), metadatumOpts).Err
			if err != nil {
				return diag.Errorf("Error adding vkcs_keymanager_secret %s metadata %s: %s", d.Id(), key, err)
			}
		}
	}

	return resourceKeyManagerSecretRead(ctx, d, meta)
}

func resourceKeyManagerSecretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	kmClient, err := config.KeyManagerV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS keymanager client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"DELETED"},
		Refresh:    keyManagerSecretWaitForSecretDeletion(kmClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
