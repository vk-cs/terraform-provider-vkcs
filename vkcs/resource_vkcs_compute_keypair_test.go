package vkcs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
)

func TestAccComputeKeypair_basic(t *testing.T) {
	var keypair keypairs.KeyPair

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeKeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeKeypairBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeKeypairExists("vkcs_compute_keypair.kp_1", &keypair),
				),
			},
		},
	})
}

func TestAccComputeKeypair_generatePrivate(t *testing.T) {
	var keypair keypairs.KeyPair

	fingerprintRe := regexp.MustCompile(`[a-f0-9:]+`)
	privateKeyRe := regexp.MustCompile(`.*BEGIN RSA PRIVATE KEY.*`)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeKeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeKeypairGeneratePrivate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeKeypairExists("vkcs_compute_keypair.kp_1", &keypair),
					resource.TestMatchResourceAttr(
						"vkcs_compute_keypair.kp_1", "fingerprint", fingerprintRe),
					resource.TestMatchResourceAttr(
						"vkcs_compute_keypair.kp_1", "private_key", privateKeyRe),
				),
			},
		},
	})
}

func testAccCheckComputeKeypairDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	computeClient, err := config.ComputeV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_compute_keypair" {
			continue
		}

		_, err := keypairs.Get(computeClient, rs.Primary.ID, keypairs.GetOpts{}).Extract()
		if err == nil {
			return fmt.Errorf("Keypair still exists")
		}
	}

	return nil
}

func testAccCheckComputeKeypairExists(n string, kp *keypairs.KeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(configer)
		computeClient, err := config.ComputeV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS compute client: %s", err)
		}

		found, err := keypairs.Get(computeClient, rs.Primary.ID, keypairs.GetOpts{}).Extract()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Keypair not found")
		}

		*kp = *found

		return nil
	}
}

const testAccComputeKeypairBasic = `
resource "vkcs_compute_keypair" "kp_1" {
  name = "kp_1"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}
`

const testAccComputeKeypairGeneratePrivate = `
resource "vkcs_compute_keypair" "kp_1" {
  name = "kp_1"
}
`
