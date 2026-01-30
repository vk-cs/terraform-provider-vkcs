package kubernetes_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesClusterV2DataSource_basic(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-cl-ds-v2-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesClusterV2Basic, map[string]string{
		"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
		"ClusterName":                     clusterName,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: clusterConfig,
				Check:  testAccCheckKubernetesClusterV2Exists("vkcs_kubernetes_cluster_v2.basic"),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterV2DataSource, map[string]string{
					"TestAccResourceClusterV2Basic": clusterConfig,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterV2DataSourceID("data.vkcs_kubernetes_cluster_v2.cluster"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "id", "vkcs_kubernetes_cluster_v2.basic", "id"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "name", "vkcs_kubernetes_cluster_v2.basic", "name"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "version", "vkcs_kubernetes_cluster_v2.basic", "version"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "cluster_type", "vkcs_kubernetes_cluster_v2.basic", "cluster_type"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "availability_zones.#", "vkcs_kubernetes_cluster_v2.basic", "availability_zones.#"),
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_cluster_v2.cluster", "availability_zones.#", "1"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "availability_zones.0", "vkcs_kubernetes_cluster_v2.basic", "availability_zones.0"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "master_count", "vkcs_kubernetes_cluster_v2.basic", "master_count"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "master_flavor", "vkcs_kubernetes_cluster_v2.basic", "master_flavor"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "network_plugin", "vkcs_kubernetes_cluster_v2.basic", "network_plugin"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "pods_ipv4_cidr", "vkcs_kubernetes_cluster_v2.basic", "pods_ipv4_cidr"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "network_id", "vkcs_kubernetes_cluster_v2.basic", "network_id"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "subnet_id", "vkcs_kubernetes_cluster_v2.basic", "subnet_id"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "loadbalancer_subnet_id", "vkcs_kubernetes_cluster_v2.basic", "loadbalancer_subnet_id"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster_v2.cluster", "public_ip", "vkcs_kubernetes_cluster_v2.basic", "public_ip"),

					resource.TestCheckResourceAttr("data.vkcs_kubernetes_cluster_v2.cluster", "master_disks.#", "3"),
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_cluster_v2.cluster", "description", ""),
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_cluster_v2.cluster", "api_lb_fip", ""),
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_cluster_v2.cluster", "insecure_registries.#", "0"),
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_cluster_v2.cluster", "labels.#", "0"),
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_cluster_v2.cluster", "loadbalancer_allowed_cidrs.#", "0"),
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_cluster_v2.cluster", "node_groups.#", "0"),

					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_v2.cluster", "uuid"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_v2.cluster", "external_network_id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_v2.cluster", "k8s_config"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_v2.cluster", "created_at"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_v2.cluster", "status"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_v2.cluster", "project_id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_v2.cluster", "api_lb_vip"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_v2.cluster", "api_address"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_cluster_v2.cluster", "region"),
				),
			},
		},
	})
}

func testAccCheckKubernetesClusterV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ct, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find cluster data source: %s", n)
		}

		if ct.Primary.ID == "" {
			return fmt.Errorf("cluster data source ID is not set")
		}

		return nil
	}
}

const testAccKubernetesClusterV2DataSource = `
{{ .TestAccResourceClusterV2Basic }}

data "vkcs_kubernetes_cluster_v2" "cluster" {
  id = vkcs_kubernetes_cluster_v2.basic.id
}
`
