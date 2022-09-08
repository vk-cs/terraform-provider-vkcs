package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatabaseCluster_importBasic(t *testing.T) {
	resourceName := "vkcs_db_cluster.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseClusterBasic,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"volume_type", "availability_zone", "disk_autoexpand", "network"},
			},
		},
	})
}
