package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatabaseClusterWithShards_basic(t *testing.T) {
	var cluster dbClusterResp

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseClusterWithShardsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseClusterWithShardsBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseClusterExists(
						"vkcs_db_cluster_with_shards.basic", &cluster),
					resource.TestCheckResourceAttrPtr(
						"vkcs_db_cluster_with_shards.basic", "name", &cluster.Name),
				),
			},
		},
	})
}

func testAccCheckDatabaseClusterWithShardsDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)

	DatabaseClient, err := config.DatabaseV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating VKCS database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_db_cluster_with_shards" {
			continue
		}
		_, err := dbClusterGet(DatabaseClient, rs.Primary.ID).extract()
		if err == nil {
			return fmt.Errorf("cluster still exists")
		}
	}

	return nil
}

var testAccDatabaseClusterWithShardsBasic = fmt.Sprintf(`
%s

%s

 resource "vkcs_db_cluster_with_shards" "basic" {
	name      = "basic"

	datastore {
	  version = "20.8"
	  type    = "clickhouse"
	}
  
  
	shard {
	  size = 1
	  shard_id = "shard0"
	  flavor_id = data.vkcs_compute_flavor.base.id
	  volume_size      = 8
	  volume_type = "ceph-ssd"
	  network {
		  uuid = vkcs_networking_network.base.id
	  }
	  availability_zone = "GZ1"
	}

	depends_on = [
		vkcs_networking_network.base,
		vkcs_networking_subnet.base
	]
 }
`, testAccBaseFlavor, testAccBaseNetwork)
