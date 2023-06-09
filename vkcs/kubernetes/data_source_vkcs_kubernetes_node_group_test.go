package kubernetes_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/nodegroups"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKubernetesNodeGroupDataSource_basic_big(t *testing.T) {
	var nodeGroup nodegroups.NodeGroup
	clusterName := "tfacc-ng-basic-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesNodeGroupClusterBase, map[string]string{"ClusterName": clusterName})
	nodeGroupConfig := acctest.AccTestRenderConfig(testAccKubernetesNodeGroupBasic, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesNodeGroupClusterBase": clusterConfig})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: nodeGroupConfig,
				Check:  testAccCheckKubernetesNodeGroupExists("vkcs_kubernetes_node_group.basic", &nodeGroup),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupDataSourceBasic, map[string]string{"TestAccKubernetesNodeGroupBasic": nodeGroupConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesNodeGroupDataSourceID("data.vkcs_kubernetes_node_group.node-group"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group.node-group", "name", "vkcs_kubernetes_node_group.basic", "name"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group.node-group", "nodes.#", "vkcs_kubernetes_node_group.basic", "node_count"),
				),
			},
		},
	})
}

func testAccKubernetesNodeGroupDataSourceID(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ct, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find node group data source: %s", resourceName)
		}

		if ct.Primary.ID == "" {
			return fmt.Errorf("node group data source ID is not set")
		}

		return nil
	}
}

const testAccKubernetesNodeGroupDataSourceBasic = `
{{ .TestAccKubernetesNodeGroupBasic }}

data "vkcs_kubernetes_node_group" "node-group" {
	uuid = vkcs_kubernetes_node_group.basic.id
}
`
