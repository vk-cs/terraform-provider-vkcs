package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ikepolicies"
)

func resourceIKEPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIKEPolicyCreate,
		ReadContext:   resourceIKEPolicyRead,
		UpdateContext: resourceIKEPolicyUpdate,
		DeleteContext: resourceIKEPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Networking client. A Networking client is needed to create a VPN service. If omitted, the `region` argument of the provider is used. Changing this creates a new service.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the policy. Changing this updates the name of the existing policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The human-readable description for the policy. Changing this updates the description of the existing policy.",
			},
			"auth_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "sha1",
				Description: "The authentication hash algorithm. Valid values are sha1, sha256, sha384, sha512. Default is sha1. Changing this updates the algorithm of the existing policy.",
			},
			"encryption_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "aes-128",
				Description: "The encryption algorithm. Valid values are 3des, aes-128, aes-192 and so on. The default value is aes-128. Changing this updates the existing policy.",
			},
			"pfs": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "group5",
				Description: "The perfect forward secrecy mode. Valid values are Group2, Group5 and Group14. Default is Group5. Changing this updates the existing policy.",
			},
			"phase1_negotiation_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "main",
				Description: "The IKE mode. A valid value is main, which is the default. Changing this updates the existing policy.",
			},
			"ike_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "v1",
				Description: "The IKE mode. A valid value is v1 or v2. Default is v1. Changing this updates the existing policy.",
			},
			"lifetime": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"units": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "The units for the lifetime of the security association. Can be either seconds or kilobytes. Default is seconds.",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The value for the lifetime of the security association. Must be a positive integer. Default is 3600.",
						},
					},
				},
				Description: "The lifetime of the security association. Consists of Unit and Value.",
			},
		},
		Description: "Manages a IKE policy resource within VKCS.",
	}
}

func resourceIKEPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	lifetime := resourceIKEPolicyLifetimeCreateOpts(d.Get("lifetime").(*schema.Set))
	authAlgorithm := resourceIKEPolicyAuthAlgorithm(d.Get("auth_algorithm").(string))
	encryptionAlgorithm := resourceIKEPolicyEncryptionAlgorithm(d.Get("encryption_algorithm").(string))
	pfs := resourceIKEPolicyPFS(d.Get("pfs").(string))
	ikeVersion := resourceIKEPolicyIKEVersion(d.Get("ike_version").(string))
	phase1NegotationMode := resourceIKEPolicyPhase1NegotiationMode(d.Get("phase1_negotiation_mode").(string))

	opts := IKEPolicyCreateOpts{
		ikepolicies.CreateOpts{
			Name:                  d.Get("name").(string),
			Description:           d.Get("description").(string),
			Lifetime:              &lifetime,
			AuthAlgorithm:         authAlgorithm,
			EncryptionAlgorithm:   encryptionAlgorithm,
			PFS:                   pfs,
			IKEVersion:            ikeVersion,
			Phase1NegotiationMode: phase1NegotationMode,
		},
	}
	log.Printf("[DEBUG] Create IKE policy: %#v", opts)

	policy, err := ikepolicies.Create(networkingClient, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForIKEPolicyCreation(networkingClient, policy.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for vkcs_vpnaas_ike_policy %s to become active: %s", policy.ID, err)
	}

	log.Printf("[DEBUG] IKE policy created: %#v", policy)

	d.SetId(policy.ID)

	return resourceIKEPolicyRead(ctx, d, meta)
}

func resourceIKEPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about IKE policy: %s", d.Id())

	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	policy, err := ikepolicies.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "IKE policy"))
	}

	log.Printf("[DEBUG] Read VKCS IKE Policy %s: %#v", d.Id(), policy)

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("auth_algorithm", policy.AuthAlgorithm)
	d.Set("encryption_algorithm", policy.EncryptionAlgorithm)
	d.Set("pfs", policy.PFS)
	d.Set("phase1_negotiation_mode", policy.Phase1NegotiationMode)
	d.Set("ike_version", policy.IKEVersion)
	d.Set("region", getRegion(d, config))

	// Set the lifetime
	lifetimeMap := make(map[string]interface{})
	lifetimeMap["units"] = policy.Lifetime.Units
	lifetimeMap["value"] = policy.Lifetime.Value
	var lifetime []map[string]interface{}
	lifetime = append(lifetime, lifetimeMap)
	if err := d.Set("lifetime", &lifetime); err != nil {
		log.Printf("[WARN] unable to set IKE policy lifetime")
	}

	return nil
}

func resourceIKEPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	opts := ikepolicies.UpdateOpts{}

	var hasChange bool

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
		hasChange = true
	}

	if d.HasChange("pfs") {
		opts.PFS = resourceIKEPolicyPFS(d.Get("pfs").(string))
		hasChange = true
	}
	if d.HasChange("auth_algorithm") {
		opts.AuthAlgorithm = resourceIKEPolicyAuthAlgorithm(d.Get("auth_algorithm").(string))
		hasChange = true
	}
	if d.HasChange("encryption_algorithm") {
		opts.EncryptionAlgorithm = resourceIKEPolicyEncryptionAlgorithm(d.Get("encryption_algorithm").(string))
		hasChange = true
	}
	if d.HasChange("phase_1_negotiation_mode") {
		opts.Phase1NegotiationMode = resourceIKEPolicyPhase1NegotiationMode(d.Get("phase_1_negotiation_mode").(string))
		hasChange = true
	}
	if d.HasChange("ike_version") {
		opts.IKEVersion = resourceIKEPolicyIKEVersion(d.Get("ike_version").(string))
		hasChange = true
	}

	if d.HasChange("lifetime") {
		lifetime := resourceIKEPolicyLifetimeUpdateOpts(d.Get("lifetime").(*schema.Set))
		opts.Lifetime = &lifetime
		hasChange = true
	}

	log.Printf("[DEBUG] Updating IKE policy with id %s: %#v", d.Id(), opts)

	if hasChange {
		err = ikepolicies.Update(networkingClient, d.Id(), opts).Err
		if err != nil {
			return diag.FromErr(err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"ACTIVE"},
			Refresh:    waitForIKEPolicyUpdate(networkingClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIKEPolicyRead(ctx, d, meta)
}

func resourceIKEPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy IKE policy: %s", d.Id())

	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForIKEPolicyDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForIKEPolicyDeletion(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := ikepolicies.Delete(networkingClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		return nil, "ACTIVE", err
	}
}

func waitForIKEPolicyCreation(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ikepolicies.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func waitForIKEPolicyUpdate(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ikepolicies.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func resourceIKEPolicyAuthAlgorithm(v string) ikepolicies.AuthAlgorithm {
	var authAlgorithm ikepolicies.AuthAlgorithm
	switch v {
	case "sha1":
		authAlgorithm = ikepolicies.AuthAlgorithmSHA1
	case "sha256":
		authAlgorithm = ikepolicies.AuthAlgorithmSHA256
	case "sha384":
		authAlgorithm = ikepolicies.AuthAlgorithmSHA384
	case "sha512":
		authAlgorithm = ikepolicies.AuthAlgorithmSHA512
	}

	return authAlgorithm
}

func resourceIKEPolicyEncryptionAlgorithm(v string) ikepolicies.EncryptionAlgorithm {
	var encryptionAlgorithm ikepolicies.EncryptionAlgorithm
	switch v {
	case "3des":
		encryptionAlgorithm = ikepolicies.EncryptionAlgorithm3DES
	case "aes-128":
		encryptionAlgorithm = ikepolicies.EncryptionAlgorithmAES128
	case "aes-192":
		encryptionAlgorithm = ikepolicies.EncryptionAlgorithmAES192
	case "aes-256":
		encryptionAlgorithm = ikepolicies.EncryptionAlgorithmAES256
	}

	return encryptionAlgorithm
}

func resourceIKEPolicyPFS(v string) ikepolicies.PFS {
	var pfs ikepolicies.PFS
	switch v {
	case "group5":
		pfs = ikepolicies.PFSGroup5
	case "group2":
		pfs = ikepolicies.PFSGroup2
	case "group14":
		pfs = ikepolicies.PFSGroup14
	}
	return pfs
}

func resourceIKEPolicyIKEVersion(v string) ikepolicies.IKEVersion {
	var ikeVersion ikepolicies.IKEVersion
	switch v {
	case "v1":
		ikeVersion = ikepolicies.IKEVersionv1
	case "v2":
		ikeVersion = ikepolicies.IKEVersionv2
	}
	return ikeVersion
}

func resourceIKEPolicyPhase1NegotiationMode(v string) ikepolicies.Phase1NegotiationMode {
	var phase1NegotiationMode ikepolicies.Phase1NegotiationMode
	if v == "main" {
		phase1NegotiationMode = ikepolicies.Phase1NegotiationModeMain
	}

	return phase1NegotiationMode
}

func resourceIKEPolicyUnit(v string) ikepolicies.Unit {
	var unit ikepolicies.Unit
	switch v {
	case "kilobytes":
		unit = ikepolicies.UnitKilobytes
	case "seconds":
		unit = ikepolicies.UnitSeconds
	}
	return unit
}

func resourceIKEPolicyLifetimeCreateOpts(d *schema.Set) ikepolicies.LifetimeCreateOpts {
	lifetimeCreateOpts := ikepolicies.LifetimeCreateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		lifetimeCreateOpts.Units = resourceIKEPolicyUnit(rawMap["units"].(string))

		value := rawMap["value"].(int)
		lifetimeCreateOpts.Value = value
	}
	return lifetimeCreateOpts
}

func resourceIKEPolicyLifetimeUpdateOpts(d *schema.Set) ikepolicies.LifetimeUpdateOpts {
	lifetimeUpdateOpts := ikepolicies.LifetimeUpdateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		lifetimeUpdateOpts.Units = resourceIKEPolicyUnit(rawMap["units"].(string))

		value := rawMap["value"].(int)
		lifetimeUpdateOpts.Value = value
	}
	return lifetimeUpdateOpts
}
