package vpnaas_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	iipsecpolicies "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/vpnaas/v2/ipsecpolicies"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func TestAccVPNaaSIPSecPolicy_basic(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
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
		ProviderFactories: acctest.AccTestProviders,
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
		ProviderFactories: acctest.AccTestProviders,
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
		ProviderFactories: acctest.AccTestProviders,
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
	config := acctest.AccTestProvider.Meta().(clients.Config)
	networkingClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.DefaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_vpnaas_ipsec_policy" {
			continue
		}
		_, err = iipsecpolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IPSec policy (%s) still exists", rs.Primary.ID)
		}
		if !errutil.Is(err, 404) {
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

		config := acctest.AccTestProvider.Meta().(clients.Config)
		networkingClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.DefaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		found, err := iipsecpolicies.Get(networkingClient, rs.Primary.ID).Extract()
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
	sdn = "neutron"
}
`

const testAccIPSecPolicyUpdate = `
resource "vkcs_vpnaas_ipsec_policy" "policy_1" {
	name = "updatedname"
	sdn = "neutron"
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
	sdn = "neutron"
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
	sdn = "neutron"
}
`
