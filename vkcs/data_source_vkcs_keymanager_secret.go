package vkcs

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
)

func getDateFilters() [4]string {
	return [4]string{
		string(secrets.DateFilterGT),
		string(secrets.DateFilterGTE),
		string(secrets.DateFilterLT),
		string(secrets.DateFilterLTE),
	}
}

func getDateFiltersRegexPreformatted() string {
	df := getDateFilters()
	return strings.Join(df[:], "|")
}

func dataSourceKeyManagerSecret() *schema.Resource {
	ret := &schema.Resource{
		ReadContext: dataSourceKeyManagerSecretRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region in which to obtain the KeyManager client. A KeyManager client is needed to fetch a secret. If omitted, the `region` argument of the provider is used.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Secret name.",
			},

			"bit_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The Secret bit length.",
			},

			"algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Secret algorithm.",
			},

			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Secret mode.",
			},

			"secret_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"symmetric", "public", "private", "passphrase", "certificate", "opaque",
				}, false),
				Description: "The Secret type. For more information see [Secret types](https://docs.openstack.org/barbican/latest/api/reference/secret_types.html).",
			},

			"acl_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Select the Secret with an ACL that contains the user. Project scope is ignored. Defaults to `false`.",
			},

			"expiration_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: dataSourceValidateDateFilter,
				Description:  "Date filter to select the Secret with expiration matching the specified criteria. See Date Filters below for more detail.",
			},

			"created_at_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: dataSourceValidateDateFilter,
				Description:  "Date filter to select the Secret with created matching the specified criteria. See Date Filters below for more detail.",
			},

			"updated_at_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: dataSourceValidateDateFilter,
				Description:  "Date filter to select the Secret with updated matching the specified criteria. See Date Filters below for more detail.",
			},

			// computed
			"acl": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of ACLs assigned to a secret.",
			},

			"secret_ref": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret reference / where to find the secret.",
			},

			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the secret.",
			},

			"expiration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the secret will expire.",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the secret was created.",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the secret was last updated.",
			},

			"content_types": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The map of the content types, assigned on the secret.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the secret.",
			},

			"payload": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The secret payload.",
			},

			"payload_content_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Secret content type.",
			},

			"payload_content_encoding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Secret encoding.",
			},

			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The map of metadata, assigned on the secret, which has been explicitly and implicitly added.",
			},
		},
		Description: "Use this data source to get the ID and the payload of an available Key secret\n\n~> **Important Security Notice** The payload of this data source will be stored *unencrypted* in your Terraform state file. **Use of this resource for production deployments is *not* recommended**. [Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).",
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

func dataSourceKeyManagerSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	kmClient, err := config.KeyManagerV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS keymanager client: %s", err)
	}

	aclOnly := d.Get("acl_only").(bool)

	listOpts := secrets.ListOpts{
		Name:            d.Get("name").(string),
		Bits:            d.Get("bit_length").(int),
		Alg:             d.Get("algorithm").(string),
		Mode:            d.Get("mode").(string),
		SecretType:      secrets.SecretType(d.Get("secret_type").(string)),
		ACLOnly:         &aclOnly,
		CreatedQuery:    dataSourceParseDateFilter(d.Get("created_at_filter").(string)),
		UpdatedQuery:    dataSourceParseDateFilter(d.Get("updated_at_filter").(string)),
		ExpirationQuery: dataSourceParseDateFilter(d.Get("expiration_filter").(string)),
	}

	log.Printf("[DEBUG] %#+v List Options: %#v", dataSourceParseDateFilter(d.Get("updated_at_filter").(string)), listOpts)

	allPages, err := secrets.List(kmClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query vkcs_keymanager_secret secrets: %s", err)
	}

	allSecrets, err := secrets.ExtractSecrets(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_keymanager_secret secrets: %s", err)
	}

	if len(allSecrets) < 1 {
		return diag.Errorf("Your query returned no vkcs_keymanager_secret results. " +
			"Please change your search criteria and try again")
	}

	if len(allSecrets) > 1 {
		log.Printf("[DEBUG] Multiple vkcs_keymanager_secret results found: %#v", allSecrets)
		return diag.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria")
	}

	secret := allSecrets[0]

	log.Printf("[DEBUG] Retrieved vkcs_keymanager_secret %s: %#v", d.Id(), secret)

	uuid := keyManagerSecretGetUUIDfromSecretRef(secret.SecretRef)

	d.SetId(uuid)
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
		log.Printf("[DEBUG] Unable to get %s secret metadata: %s", uuid, err)
	}
	d.Set("metadata", metadataMap)

	if secret.Expiration == (time.Time{}) {
		d.Set("expiration", "")
	} else {
		d.Set("expiration", secret.Expiration.Format(time.RFC3339))
	}

	acl, err := acls.GetSecretACL(kmClient, d.Id()).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get %s secret acls: %s", uuid, err)
	}
	d.Set("acl", flattenKeyManagerACLs(acl))

	// Set the region
	d.Set("region", getRegion(d, config))

	return nil
}

func dataSourceParseDateFilter(date string) *secrets.DateQuery {
	// error checks are not necessary, since they were validated by terraform validate functions
	var parts []string
	if regexp.MustCompile("^" + getDateFiltersRegexPreformatted() + ":").Match([]byte(date)) {
		parts = strings.SplitN(date, ":", 2)
	} else {
		parts = []string{date}
	}

	var parsedTime time.Time
	var filter *secrets.DateQuery

	if len(parts) == 2 {
		parsedTime, _ = time.Parse(time.RFC3339, parts[1])

		filter = &secrets.DateQuery{Date: parsedTime, Filter: secrets.DateFilter(parts[0])}
	} else {
		parsedTime, _ = time.Parse(time.RFC3339, parts[0])

		filter = &secrets.DateQuery{Date: parsedTime}
	}

	if parsedTime == (time.Time{}) {
		return nil
	}

	return filter
}

func dataSourceValidateDateFilter(v interface{}, k string) (ws []string, errors []error) {
	var parts []string
	if regexp.MustCompile("^" + getDateFiltersRegexPreformatted() + ":").Match([]byte(v.(string))) {
		parts = strings.SplitN(v.(string), ":", 2)
	} else {
		parts = []string{v.(string)}
	}

	if len(parts) == 2 {
		supportedDateFilters := getDateFilters()
		if !strSliceContains(supportedDateFilters[:], parts[0]) {
			errors = append(errors, fmt.Errorf("invalid %q date filter, supported: %+q", parts[0], supportedDateFilters))
		}

		_, err := time.Parse(time.RFC3339, parts[1])
		if err != nil {
			errors = append(errors, err)
		}

		return
	}

	_, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		errors = append(errors, err)
	}

	return
}
