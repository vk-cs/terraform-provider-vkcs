package kubernetes_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesClusterAddonV2DataSource_basic(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-clad-ds-v2-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesClusterV2Base, map[string]string{
		"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
		"ClusterName":                     clusterName,
	})
	clusterWithNodeGroupConfig := acctest.AccTestRenderConfig(testAccKubernetesNodeGroupV2Base, map[string]string{
		"TestAccKubernetesNodeGroupV2ClusterBase": clusterConfig,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterAddonV2DataSourceResource, map[string]string{
					"TestAccKubernetesClusterAddonV2DataSourceClusterBase": clusterWithNodeGroupConfig,
				}),
				Check: testAccCheckKubernetesClusterAddonV2DataSourceResourceExists("vkcs_kubernetes_cluster_addon_v2.ingress-nginx"),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterAddonV2DataSourceConfig, map[string]string{
					"TestAccKubernetesClusterAddonV2DataSourceResource": acctest.AccTestRenderConfig(testAccKubernetesClusterAddonV2DataSourceResource, map[string]string{
						"TestAccKubernetesClusterAddonV2DataSourceClusterBase": clusterConfig,
					}),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterAddonV2DataSourceID("data.vkcs_kubernetes_cluster_addon_v2.ingress-nginx"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "id", "vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "id"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "cluster_id", "vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "cluster_id"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "addon_id", "vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "addon_id"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "addon_name", "vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "addon_name"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "addon_version_id", "vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "addon_version_id"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "namespace", "vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "namespace"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "values", "vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "values"),

					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "created_at"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "updated_at"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "status"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_addon_v2.ingress-nginx", "region"),
				),
			},
		},
	})
}

func testAccCheckKubernetesClusterAddonV2DataSourceResourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Kubernetes cluster addon resource not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cluster addon resource ID is not set")
		}
		return nil
	}
}

func testAccCheckKubernetesClusterAddonV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find cluster addon data source: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cluster addon data source ID is not set")
		}
		return nil
	}
}

const testAccKubernetesClusterAddonV2DataSourceResource = `
{{ .TestAccKubernetesClusterAddonV2DataSourceClusterBase }}

data "vkcs_kubernetes_addon_v2" "ingress-nginx" {
  name    = "ingress-nginx"
  version = "4.12.1"
}

resource "vkcs_kubernetes_cluster_addon_v2" "ingress-nginx" {
  cluster_id       = vkcs_kubernetes_cluster_v2.base.id
  addon_id         = data.vkcs_kubernetes_addon_v2.ingress-nginx.addon_id
  addon_name       = "ingress-nginx"
  addon_version_id = data.vkcs_kubernetes_addon_v2.ingress-nginx.id
  namespace        = "ingress-nginx"
  values           = data.vkcs_kubernetes_addon_v2.ingress-nginx.values_template
}
`

const testAccKubernetesClusterAddonV2DataSourceConfig = `
{{ .TestAccKubernetesClusterAddonV2DataSourceResource }}

data "vkcs_kubernetes_cluster_addon_v2" "ingress-nginx" {
  id = vkcs_kubernetes_cluster_addon_v2.ingress-nginx.id
}
`
