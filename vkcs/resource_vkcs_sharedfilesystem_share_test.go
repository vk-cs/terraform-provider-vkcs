package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

func TestAccSFSShare_basic(t *testing.T) {
	var share shares.Share

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckSFS(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "description", "test share description"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "share_proto", "NFS"),
				),
			},
			{
				Config: testAccSFSShareConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share_updated"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "description", ""),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "share_proto", "NFS"),
				),
			},
			{
				Config: testAccSFSShareConfigExtend,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share_extended"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "share_proto", "NFS"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "size", "2"),
				),
			},
		},
	})
}

func TestAccSFSShare_update(t *testing.T) {
	var share shares.Share

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckSFS(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareConfigMetadataUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "description", "test share description"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "share_proto", "NFS"),
					testAccCheckSFSShareMetadataEquals("key", "value", &share),
				),
			},
			{
				Config: testAccSFSShareConfigMetadataUpdate1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share"),
					testAccCheckSFSShareMetadataEquals("key", "value", &share),
					testAccCheckSFSShareMetadataEquals("new_key", "new_value", &share),
				),
			},
			{
				Config: testAccSFSShareConfigMetadataUpdate2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share"),
					testAccCheckSFSShareMetadataAbsent("key", &share),
					testAccCheckSFSShareMetadataEquals("new_key", "new_value", &share),
				),
			},
			{
				Config: testAccSFSShareConfigMetadataUpdate3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share"),
					testAccCheckSFSShareMetadataAbsent("key", &share),
					testAccCheckSFSShareMetadataAbsent("new_key", &share),
				),
			},
		},
	})
}

func TestAccSFSShare_admin(t *testing.T) {
	var share shares.Share

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckSFS(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareAdminConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share_admin"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "description", "test share description"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "share_proto", "NFS"),
				),
			},
			{
				Config: testAccSFSShareAdminConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSShareExists("vkcs_sharedfilesystem_share.share_1", &share),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "name", "nfs_share_admin_updated"),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "description", ""),
					resource.TestCheckResourceAttr("vkcs_sharedfilesystem_share.share_1", "share_proto", "NFS"),
				),
			},
		},
	})
}

func testAccCheckSFSShareDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config)
	sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_sharedfilesystem_securityservice" {
			continue
		}

		_, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Manila share still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckSFSShareExists(n string, share *shares.Share) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
		}

		found, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Share not found")
		}

		*share = *found

		return nil
	}
}

func testAccCheckSFSShareMetadataEquals(key string, value string, share *shares.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config)
		sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
		}

		metadatum, err := shares.GetMetadatum(sfsClient, share.ID, key).Extract()
		if err != nil {
			return err
		}

		if metadatum[key] != value {
			return fmt.Errorf("Metadata does not match. Expected %v but got %v", metadatum, value)
		}

		return nil
	}
}

func testAccCheckSFSShareMetadataAbsent(key string, share *shares.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config)
		sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
		}

		_, err = shares.GetMetadatum(sfsClient, share.ID, key).Extract()
		if err == nil {
			return fmt.Errorf("Metadata %s key must not exist", key)
		}

		return nil
	}
}

const testAccSFSShareConfigBasic = `
resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`

const testAccSFSShareConfigUpdate = `
resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share_updated"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`

const testAccSFSShareConfigExtend = `
resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share_extended"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 2
}
`

//const testAccSFSShareConfigShrink = `
//resource "vkcs_sharedfilesystem_share" "share_1" {
//  name             = "nfs_share_shrunk"
//  share_proto      = "NFS"
//  share_type       = "dhss_false"
//  size             = 1
//}
//`

const testAccSFSShareConfigMetadataUpdate = `
resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`

const testAccSFSShareConfigMetadataUpdate1 = `
resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`

const testAccSFSShareConfigMetadataUpdate2 = `
resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`

const testAccSFSShareConfigMetadataUpdate3 = `
resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`

const testAccSFSShareAdminConfigBasic = `
resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share_admin"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`

const testAccSFSShareAdminConfigUpdate = `
resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share_admin_updated"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`
