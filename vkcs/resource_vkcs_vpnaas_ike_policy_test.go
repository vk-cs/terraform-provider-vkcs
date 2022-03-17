package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ikepolicies"
)

func TestAccVPNaaSIKEPolicy_basic(t *testing.T) {
	var policy ikepolicies.Policy
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
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
		PreCheck: func() {
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
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
		PreCheck: func() {
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
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
		PreCheck: func() {
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
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
	config := testAccProvider.Meta().(*config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_vpnaas_ike_policy" {
			continue
		}
		_, err = ikepolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IKE policy (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
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

		config := testAccProvider.Meta().(*config)
		networkingClient, err := config.NetworkingV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		found, err := ikepolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*policy = *found

		return nil
	}
}

const testAccIKEPolicyBasic = `
resource "vkcs_vpnaas_ike_policy" "policy_1" {
}
`

const testAccIKEPolicyUpdate = `
resource "vkcs_vpnaas_ike_policy" "policy_1" {
	name = "updatedname"
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
