package kubernetes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesAddonDataSource_basic_big(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesAddonDataSourceBasic, map[string]string{
					"TestAccKubernetesAddonNetworkingBase": testAccKubernetesAddonNetworkingBase,
					"TestAccKubernetesAddonClusterBase":    testAccKubernetesAddonClusterBase,
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
{{ .TestAccKubernetesAddonNetworkingBase }}
{{ .TestAccKubernetesAddonClusterBase }}

data "vkcs_kubernetes_addon" "ingress-nginx" {
  cluster_id = vkcs_kubernetes_cluster.cluster.id
  name       = "ingress-nginx"
  version    = "4.1.4"
}
`
