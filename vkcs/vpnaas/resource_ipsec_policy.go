package vpnaas

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/vpnaas"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/networking"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
)

func ResourceIPSecPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPSecPolicyCreate,
		ReadContext:   resourceIPSecPolicyRead,
		UpdateContext: resourceIPSecPolicyUpdate,
		DeleteContext: resourceIPSecPolicyDelete,
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
				Description: "The region in which to obtain the Networking client. A Networking client is needed to create an IPSec policy. If omitted, the `region` argument of the provider is used. Changing this creates a new policy.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the policy. Changing this updates the name of the existing policy.",
			},
			"auth_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The authentication hash algorithm. Valid values are sha1, sha256, sha384, sha512. Default is sha1. Changing this updates the algorithm of the existing policy.",
			},
			"encapsulation_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The encapsulation mode. Valid values are tunnel and transport. Default is tunnel. Changing this updates the existing policy.",
			},
			"pfs": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The perfect forward secrecy mode. Valid values are Group2, Group5 and Group14. Default is Group5. Changing this updates the existing policy.",
			},
			"encryption_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The encryption algorithm. Valid values are 3des, aes-128, aes-192 and so on. The default value is aes-128. Changing this updates the existing policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The human-readable description for the policy. Changing this updates the description of the existing policy.",
			},
			"transform_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The transform protocol. Valid values are ESP, AH and AH-ESP. Changing this updates the existing policy. Default is ESP.",
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
			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateDiagFunc: networking.ValidateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is project's default SDN.",
			},
		},
		Description: "Manages a IPSec policy resource within VKCS.",
	}
}

func resourceIPSecPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	encapsulationMode := resourceIPSecPolicyEncapsulationMode(d.Get("encapsulation_mode").(string))
	authAlgorithm := resourceIPSecPolicyAuthAlgorithm(d.Get("auth_algorithm").(string))
	encryptionAlgorithm := resourceIPSecPolicyEncryptionAlgorithm(d.Get("encryption_algorithm").(string))
	pfs := resourceIPSecPolicyPFS(d.Get("pfs").(string))
	transformProtocol := resourceIPSecPolicyTransformProtocol(d.Get("transform_protocol").(string))
	lifetime := resourceIPSecPolicyLifetimeCreateOpts(d.Get("lifetime").(*schema.Set))

	opts := IPSecPolicyCreateOpts{
		CreateOpts: ipsecpolicies.CreateOpts{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			EncapsulationMode:   encapsulationMode,
			AuthAlgorithm:       authAlgorithm,
			EncryptionAlgorithm: encryptionAlgorithm,
			PFS:                 pfs,
			TransformProtocol:   transformProtocol,
			Lifetime:            &lifetime,
		},
	}

	log.Printf("[DEBUG] Create IPSec policy: %#v", opts)

	policy, err := ipsecpolicies.Create(networkingClient, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(policy.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForIPSecPolicyCreation(networkingClient, policy.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for vkcs_vpnaas_ipsec_policy %s to become active: %s", policy.ID, err)
	}

	log.Printf("[DEBUG] IPSec policy created: %#v", policy)

	return resourceIPSecPolicyRead(ctx, d, meta)
}

func resourceIPSecPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about IPSec policy: %s", d.Id())

	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var policy ipsecPolicyExtended
	err = vpnaas.ExtractIPSecPolicyInto(ipsecpolicies.Get(networkingClient, d.Id()), &policy)
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "IPSec policy"))
	}

	log.Printf("[DEBUG] Read VKCS IPSec policy %s: %#v", d.Id(), policy)

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("encapsulation_mode", policy.EncapsulationMode)
	d.Set("encryption_algorithm", policy.EncryptionAlgorithm)
	d.Set("transform_protocol", policy.TransformProtocol)
	d.Set("pfs", policy.PFS)
	d.Set("auth_algorithm", policy.AuthAlgorithm)
	d.Set("region", util.GetRegion(d, config))
	d.Set("sdn", policy.SDN)

	// Set the lifetime
	lifetimeMap := make(map[string]interface{})
	lifetimeMap["units"] = policy.Lifetime.Units
	lifetimeMap["value"] = policy.Lifetime.Value
	var lifetime []map[string]interface{}
	lifetime = append(lifetime, lifetimeMap)
	if err := d.Set("lifetime", &lifetime); err != nil {
		log.Printf("[WARN] unable to set IPSec policy lifetime")
	}

	return nil
}

func resourceIPSecPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var hasChange bool
	opts := ipsecpolicies.UpdateOpts{}

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

	if d.HasChange("auth_algorithm") {
		opts.AuthAlgorithm = resourceIPSecPolicyAuthAlgorithm(d.Get("auth_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("encryption_algorithm") {
		opts.EncryptionAlgorithm = resourceIPSecPolicyEncryptionAlgorithm(d.Get("encryption_algorithm").(string))
		hasChange = true
	}

	if d.HasChange("transform_protocol") {
		opts.TransformProtocol = resourceIPSecPolicyTransformProtocol(d.Get("transform_protocol").(string))
		hasChange = true
	}

	if d.HasChange("pfs") {
		opts.PFS = resourceIPSecPolicyPFS(d.Get("pfs").(string))
		hasChange = true
	}

	if d.HasChange("encapsulation_mode") {
		opts.EncapsulationMode = resourceIPSecPolicyEncapsulationMode(d.Get("encapsulation_mode").(string))
		hasChange = true
	}

	if d.HasChange("lifetime") {
		lifetime := resourceIPSecPolicyLifetimeUpdateOpts(d.Get("lifetime").(*schema.Set))
		opts.Lifetime = &lifetime
		hasChange = true
	}

	log.Printf("[DEBUG] Updating IPSec policy with id %s: %#v", d.Id(), opts)

	if hasChange {
		_, err = ipsecpolicies.Update(networkingClient, d.Id(), opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"ACTIVE"},
			Refresh:    waitForIPSecPolicyUpdate(networkingClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceIPSecPolicyRead(ctx, d, meta)
}

func resourceIPSecPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy IPSec policy: %s", d.Id())

	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForIPSecPolicyDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForIPSecPolicyDeletion(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := ipsecpolicies.Delete(networkingClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if _, ok := err.(gophercloud.ErrDefault409); ok {
			return nil, "ACTIVE", nil
		}

		return nil, "ACTIVE", err
	}
}

func waitForIPSecPolicyCreation(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ipsecpolicies.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func waitForIPSecPolicyUpdate(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := ipsecpolicies.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return policy, "ACTIVE", nil
	}
}

func resourceIPSecPolicyTransformProtocol(trp string) ipsecpolicies.TransformProtocol {
	var protocol ipsecpolicies.TransformProtocol
	switch trp {
	case "esp":
		protocol = ipsecpolicies.TransformProtocolESP
	case "ah":
		protocol = ipsecpolicies.TransformProtocolAH
	case "ah-esp":
		protocol = ipsecpolicies.TransformProtocolAHESP
	}
	return protocol
}
func resourceIPSecPolicyPFS(pfsString string) ipsecpolicies.PFS {
	var pfs ipsecpolicies.PFS
	switch pfsString {
	case "group2":
		pfs = ipsecpolicies.PFSGroup2
	case "group5":
		pfs = ipsecpolicies.PFSGroup5
	case "group14":
		pfs = ipsecpolicies.PFSGroup14
	}
	return pfs
}
func resourceIPSecPolicyEncryptionAlgorithm(encryptionAlgo string) ipsecpolicies.EncryptionAlgorithm {
	var alg ipsecpolicies.EncryptionAlgorithm
	switch encryptionAlgo {
	case "3des":
		alg = ipsecpolicies.EncryptionAlgorithm3DES
	case "aes-128":
		alg = ipsecpolicies.EncryptionAlgorithmAES128
	case "aes-256":
		alg = ipsecpolicies.EncryptionAlgorithmAES256
	case "aes-192":
		alg = ipsecpolicies.EncryptionAlgorithmAES192
	}
	return alg
}
func resourceIPSecPolicyAuthAlgorithm(authAlgo string) ipsecpolicies.AuthAlgorithm {
	var alg ipsecpolicies.AuthAlgorithm
	switch authAlgo {
	case "sha1":
		alg = ipsecpolicies.AuthAlgorithmSHA1
	case "sha256":
		alg = ipsecpolicies.AuthAlgorithmSHA256
	case "sha384":
		alg = ipsecpolicies.AuthAlgorithmSHA384
	case "sha512":
		alg = ipsecpolicies.AuthAlgorithmSHA512
	}
	return alg
}
func resourceIPSecPolicyEncapsulationMode(encMode string) ipsecpolicies.EncapsulationMode {
	var mode ipsecpolicies.EncapsulationMode
	switch encMode {
	case "tunnel":
		mode = ipsecpolicies.EncapsulationModeTunnel
	case "transport":
		mode = ipsecpolicies.EncapsulationModeTransport
	}
	return mode
}

func resourceIPSecPolicyLifetimeCreateOpts(d *schema.Set) ipsecpolicies.LifetimeCreateOpts {
	lifetime := ipsecpolicies.LifetimeCreateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		lifetime.Units = resourceIPSecPolicyUnit(rawMap["units"].(string))

		value := rawMap["value"].(int)
		lifetime.Value = value
	}
	return lifetime
}

func resourceIPSecPolicyUnit(units string) ipsecpolicies.Unit {
	var unit ipsecpolicies.Unit
	switch units {
	case "seconds":
		unit = ipsecpolicies.UnitSeconds
	case "kilobytes":
		unit = ipsecpolicies.UnitKilobytes
	}
	return unit
}

func resourceIPSecPolicyLifetimeUpdateOpts(d *schema.Set) ipsecpolicies.LifetimeUpdateOpts {
	lifetimeUpdateOpts := ipsecpolicies.LifetimeUpdateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		lifetimeUpdateOpts.Units = resourceIPSecPolicyUnit(rawMap["units"].(string))

		value := rawMap["value"].(int)
		lifetimeUpdateOpts.Value = value
	}
	return lifetimeUpdateOpts
}
