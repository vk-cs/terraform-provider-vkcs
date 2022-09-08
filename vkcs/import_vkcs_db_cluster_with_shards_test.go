package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatabaseClusterWithShards_importBasic(t *testing.T) {
	resourceName := "vkcs_db_cluster_with_shards.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseClusterWithShardsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccDatabaseClusterWithShardsBasic, testAccValues),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"shard.0.volume_type", "shard.1.volume_type", "shard.0.availability_zone", "shard.1.availability_zone", "shard.0.network", "shard.1.network", "shard.0.shard_id", "shard.0.size", "shard.1.shard_id", "shard.1.size"},
			},
		},
	})
}
