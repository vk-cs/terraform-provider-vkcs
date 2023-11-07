package vpnaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	iikepolicies "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/vpnaas/v2/ikepolicies"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ikepolicies"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func TestAccVPNaaSIKEPolicy_basic(t *testing.T) {
	var policy ikepolicies.Policy
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckIKEPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyExists("vkcs_vpnaas_ike_policy.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ike_policy.policy_1", "name", &policy.Name),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ike_policy.policy_1", "description", &policy.Description),
				),
			},
		},
	})
}

func TestAccVPNaaSIKEPolicy_withLifetime(t *testing.T) {
	var policy ikepolicies.Policy
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckIKEPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyWithLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyExists("vkcs_vpnaas_ike_policy.policy_1", &policy),
				),
			},
		},
	})
}

func TestAccVPNaaSIKEPolicy_Update(t *testing.T) {
	var policy ikepolicies.Policy
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckIKEPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyExists("vkcs_vpnaas_ike_policy.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ike_policy.policy_1", "name", &policy.Name),
				),
			},
			{
				Config: testAccIKEPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyExists("vkcs_vpnaas_ike_policy.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ike_policy.policy_1", "name", &policy.Name),
				),
			},
		},
	})
}

func TestAccVPNaaSIKEPolicy_withLifetimeUpdate(t *testing.T) {
	var policy ikepolicies.Policy
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckIKEPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyWithLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyExists("vkcs_vpnaas_ike_policy.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ike_policy.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_ike_policy.policy_1", "pfs", &policy.PFS),
				),
			},
			{
				Config: testAccIKEPolicyWithLifetimeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyExists("vkcs_vpnaas_ike_policy.policy_1", &policy),
				),
			},
		},
	})
}

func testAccCheckIKEPolicyDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	networkingClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.DefaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_vpnaas_ike_policy" {
			continue
		}
		_, err = iikepolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IKE policy (%s) still exists", rs.Primary.ID)
		}
		if !errutil.Is(err, 404) {
			return err
		}
	}
	return nil
}

func testAccCheckIKEPolicyExists(n string, policy *ikepolicies.Policy) resource.TestCheckFunc {
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

		found, err := iikepolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*policy = *found

		return nil
	}
}

const testAccIKEPolicyBasic = `
resource "vkcs_vpnaas_ike_policy" "policy_1" {
	sdn = "neutron"
}
`

const testAccIKEPolicyUpdate = `
resource "vkcs_vpnaas_ike_policy" "policy_1" {
	name = "updatedname"
	sdn = "neutron"
}
`

const testAccIKEPolicyWithLifetime = `
resource "vkcs_vpnaas_ike_policy" "policy_1" {
	auth_algorithm = "sha256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1200
	}
	sdn = "neutron"
}
`

const testAccIKEPolicyWithLifetimeUpdate = `
resource "vkcs_vpnaas_ike_policy" "policy_1" {
	auth_algorithm = "sha256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1400
	}
}
`
