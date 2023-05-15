package db_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDatabaseClusterWithShards_importBasic_big(t *testing.T) {
	resourceName := "vkcs_db_cluster_with_shards.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckDatabaseClusterWithShardsDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDatabaseClusterWithShardsBasic),
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
