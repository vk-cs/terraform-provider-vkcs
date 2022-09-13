package vkcs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

func TestAccSFSShareAccess_basic(t *testing.T) {
	var shareAccess1 shares.AccessRight
	var shareAccess2 shares.AccessRight

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSShareAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccSFSShareAccessConfigBasic, map[string]string{"TestAccSFSShareNetworkConfigBasic": testAccRenderConfig(testAccSFSShareNetworkConfigBasic, map[string]string{"TestAccSFSShareNetworkConfig": testAccSFSShareNetworkConfig}), "TestAccSFSShareAccessConfig": testAccSFSShareAccessConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareAccessExists("vkcs_sharedfilesystem_share_access.share_access_1", &shareAccess1),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share_access.share_access_1", "access_type", "ip"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share_access.share_access_1", "access_to", "192.168.199.10"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share_access.share_access_1", "access_level", "rw"),
					resource.TestMatchResourceAttr("vkcs_sharedfilesystem_share_access.share_access_1", "share_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					testAccCheckSFSShareAccessExists("vkcs_sharedfilesystem_share_access.share_access_2", &shareAccess2),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share_access.share_access_2", "access_type", "ip"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share_access.share_access_2", "access_to", "192.168.199.11"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share_access.share_access_2", "access_level", "rw"),
					resource.TestMatchResourceAttr("vkcs_sharedfilesystem_share_access.share_access_2", "share_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					testAccCheckSFSShareAccessDiffers(&shareAccess1, &shareAccess2),
				),
			},
			{
				Config: testAccRenderConfig(testAccSFSShareAccessConfigUpdate, map[string]string{"TestAccSFSShareNetworkConfigBasic": testAccRenderConfig(testAccSFSShareNetworkConfigBasic, map[string]string{"TestAccSFSShareNetworkConfig": testAccSFSShareNetworkConfig}), "TestAccSFSShareAccessConfig": testAccSFSShareAccessConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareAccessExists("vkcs_sharedfilesystem_share_access.share_access_1", &shareAccess1),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share_access.share_access_1", "access_type", "ip"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share_access.share_access_1", "access_to", "192.168.199.10"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share_access.share_access_1", "access_level", "ro"),
					resource.TestMatchResourceAttr("vkcs_sharedfilesystem_share_access.share_access_1", "share_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					testAccCheckSFSShareAccessExists("vkcs_sharedfilesystem_share_access.share_access_2", &shareAccess2),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share_access.share_access_2", "access_type", "ip"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share_access.share_access_2", "access_to", "192.168.199.11"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share_access.share_access_2", "access_level", "ro"),
					resource.TestMatchResourceAttr("vkcs_sharedfilesystem_share_access.share_access_2", "share_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					testAccCheckSFSShareAccessDiffers(&shareAccess1, &shareAccess2),
				),
			},
		},
	})
}

func testAccCheckSFSShareAccessDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config)
	sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_sharedfilesystem_share_access" {
			continue
		}

		var shareID string
		for k, v := range rs.Primary.Attributes {
			if k == "share_id" {
				shareID = v
				break
			}
		}

		access, err := shares.ListAccessRights(sfsClient, shareID).Extract()
		if err == nil {
			for _, v := range access {
				if v.ID == rs.Primary.ID {
					return fmt.Errorf("Manila share access still exists: %s", rs.Primary.ID)
				}
			}
		}
	}

	return nil
}

func testAccCheckSFSShareAccessExists(n string, share *shares.AccessRight) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config)
		sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
		}

		var shareID string
		for k, v := range rs.Primary.Attributes {
			if k == "share_id" {
				shareID = v
				break
			}
		}

		sfsClient.Microversion = sharedFilesystemMinMicroversion

		access, err := shares.ListAccessRights(sfsClient, shareID).Extract()
		if err != nil {
			return fmt.Errorf("Unable to get %s share: %s", shareID, err)
		}

		var found shares.AccessRight

		for _, v := range access {
			if v.ID == rs.Primary.ID {
				found = v
				break
			}
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("ShareAccess not found")
		}

		*share = found

		return nil
	}
}

func testAccCheckSFSShareAccessDiffers(shareAccess1, shareAccess2 *shares.AccessRight) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if shareAccess1.ID != shareAccess2.ID {
			return nil
		}
		return fmt.Errorf("Share accesses should differ")
	}
}

const testAccSFSShareAccessConfig = `
resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "default_share_type"
  size             = 1
  share_network_id = vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id
}
`

const testAccSFSShareAccessConfigBasic = `
{{.TestAccSFSShareNetworkConfigBasic}}

{{.TestAccSFSShareAccessConfig}}

resource "vkcs_sharedfilesystem_share_access" "share_access_1" {
  share_id     = vkcs_sharedfilesystem_share.share_1.id
  access_type  = "ip"
  access_to    = "192.168.199.10"
  access_level = "rw"
}

resource "vkcs_sharedfilesystem_share_access" "share_access_2" {
  share_id     = vkcs_sharedfilesystem_share.share_1.id
  access_type  = "ip"
  access_to    = "192.168.199.11"
  access_level = "rw"
}
`

const testAccSFSShareAccessConfigUpdate = `
{{.TestAccSFSShareNetworkConfigBasic}}

{{.TestAccSFSShareAccessConfig}}

resource "vkcs_sharedfilesystem_share_access" "share_access_1" {
  share_id     = vkcs_sharedfilesystem_share.share_1.id
  access_type  = "ip"
  access_to    = "192.168.199.10"
  access_level = "ro"
}

resource "vkcs_sharedfilesystem_share_access" "share_access_2" {
  share_id     = vkcs_sharedfilesystem_share.share_1.id
  access_type  = "ip"
  access_to    = "192.168.199.11"
  access_level = "ro"
}
`
