package kubernetes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/clusters"
)

func TestAccKubernetesAddonDataSource_basic_big(t *testing.T) {
	var cluster clusters.Cluster
	uniqueSuffix := acctest.GenerateNameSuffix()
	baseConfig := acctest.AccTestRenderConfig(testAccKubernetesAddonClusterBase,
		map[string]string{"TestAccKubernetesAddonNetworkingBase": testAccKubernetesAddonNetworkingBase,
			"Suffix":        uniqueSuffix,
			"NodeGroupName": uniqueKubernetesNodeGroupName(),
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: baseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterExists("vkcs_kubernetes_cluster.cluster", &cluster),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesAddonDataSourceBasic, map[string]string{
					"TestAccKubernetesAddonClusterBase": baseConfig,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_addon.ingress-nginx", "cluster_id", "vkcs_kubernetes_cluster.cluster", "id"),
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_addon.ingress-nginx", "name", "ingress-nginx"),
				),
			},
		},
	})
}

const testAccKubernetesAddonDataSourceBasic = `
{{ .TestAccKubernetesAddonClusterBase }}

data "vkcs_kubernetes_addon" "ingress-nginx" {
  cluster_id = vkcs_kubernetes_cluster.cluster.id
  name       = "ingress-nginx"
  version    = "4.7.1"
}
`
