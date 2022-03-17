package vkcs

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
)

func TestAccVPNaaSIPSecPolicy_basic(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPSecPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyExists("vkcs_vpnaas_ipsec_policy.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ipsec_policy.policy_1", "name", &policy.Name),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ipsec_policy.policy_1", "description", &policy.Description),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ipsec_policy.policy_1", "pfs", &policy.PFS),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ipsec_policy.policy_1", "transform_protocol", &policy.TransformProtocol),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ipsec_policy.policy_1", "encapsulation_mode", &policy.EncapsulationMode),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ipsec_policy.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ipsec_policy.policy_1", "encryption_algorithm", &policy.EncryptionAlgorithm),
				),
			},
		},
	})
}

func TestAccVPNaaSIPSecPolicy_withLifetime(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPSecPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyWithLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyExists("vkcs_vpnaas_ipsec_policy.policy_1", &policy),
					testAccCheckLifetime("vkcs_vpnaas_ipsec_policy.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
				),
			},
		},
	})
}

func TestAccVPNaaSIPSecPolicy_Update(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPSecPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyExists("vkcs_vpnaas_ipsec_policy.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ipsec_policy.policy_1", "name", &policy.Name),
				),
			},
			{
				Config: testAccIPSecPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyExists("vkcs_vpnaas_ipsec_policy.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ipsec_policy.policy_1", "name", &policy.Name),
				),
			},
		},
	})
}

func TestAccVPNaaSIPSecPolicy_withLifetimeUpdate(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPSecPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyWithLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyExists("vkcs_vpnaas_ipsec_policy.policy_1", &policy),
					testAccCheckLifetime("vkcs_vpnaas_ipsec_policy.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ipsec_policy.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ipsec_policy.policy_1", "pfs", &policy.PFS),
				),
			},
			{
				Config: testAccIPSecPolicyWithLifetimeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyExists("vkcs_vpnaas_ipsec_policy.policy_1", &policy),
					testAccCheckLifetime("vkcs_vpnaas_ipsec_policy.policy_1", &policy.Lifetime.Units, &policy.Lifetime.Value),
				),
			},
		},
	})
}

func testAccCheckIPSecPolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_vpnaas_ipsec_policy" {
			continue
		}
		_, err = ipsecpolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IPSec policy (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckIPSecPolicyExists(n string, policy *ipsecpolicies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config)
		networkingClient, err := config.NetworkingV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		found, err := ipsecpolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*policy = *found

		return nil
	}
}

func testAccCheckLifetime(n string, unit *string, value *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		for k, v := range rs.Primary.Attributes {
			println("[DEBUG] key:", k, "value:", v)
			if strings.HasPrefix(k, "lifetime.") && k[9] >= '0' && k[9] <= '9' && strings.HasSuffix(k, ".units") {
				index := strings.LastIndex(k, ".")
				base := k[:index]
				expectedValue := rs.Primary.Attributes[base+".value"]
				expectedUnit := rs.Primary.Attributes[k]
				println("[DEBUG] expectedValue:", expectedValue, "expectedUnit:", expectedUnit)

				if expectedUnit != *unit {
					return fmt.Errorf("Expected lifetime unit %v but found %v", expectedUnit, *unit)
				}
				if expectedValue != strconv.Itoa(*value) {
					return fmt.Errorf("Expected lifetime value %v but found %v", expectedValue, *value)
				}
			}
		}

		return nil
	}
}

const testAccIPSecPolicyBasic = `
resource "vkcs_vpnaas_ipsec_policy" "policy_1" {
}
`

const testAccIPSecPolicyUpdate = `
resource "vkcs_vpnaas_ipsec_policy" "policy_1" {
	name = "updatedname"
}
`

const testAccIPSecPolicyWithLifetime = `
resource "vkcs_vpnaas_ipsec_policy" "policy_1" {
	auth_algorithm = "sha256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1200
	}
}
`

const testAccIPSecPolicyWithLifetimeUpdate = `
resource "vkcs_vpnaas_ipsec_policy" "policy_1" {
	auth_algorithm = "sha256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1400
	}
}
`
