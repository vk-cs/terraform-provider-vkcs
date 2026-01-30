package kubernetes_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesNodeGroupV2DataSource_basic(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-ng-ds-v2-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesNodeGroupV2DataSourceClusterBase, map[string]string{
		"ClusterName":                     clusterName,
		"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
	})
	nodeGroupConfig := acctest.AccTestRenderConfig(testAccKubernetesNodeGroupV2DataSource, map[string]string{
		"TestAccKubernetesNodeGroupV2DataSourceClusterBase": clusterConfig,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: nodeGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterV2Exists("vkcs_kubernetes_cluster_v2.base"),
					testAccCheckKubernetesNodeGroupV2Exists("vkcs_kubernetes_node_group_v2.base"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.base", "name", "tfacc-ds-full-v2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.base", "scale_type", "auto_scale"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.base", "availability_zone", "MS1"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupV2DataSourceFullRead, map[string]string{
					"TestAccKubernetesNodeGroupV2DataSourceFull": nodeGroupConfig,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesNodeGroupV2DataSourceID("data.vkcs_kubernetes_node_group_v2.node_group"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "id", "vkcs_kubernetes_node_group_v2.base", "id"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "name", "vkcs_kubernetes_node_group_v2.base", "name"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "cluster_id", "vkcs_kubernetes_node_group_v2.base", "cluster_id"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "node_flavor", "vkcs_kubernetes_node_group_v2.base", "node_flavor"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "availability_zone", "vkcs_kubernetes_node_group_v2.base", "availability_zone"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "disk_type", "vkcs_kubernetes_node_group_v2.base", "disk_type"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "disk_size", "vkcs_kubernetes_node_group_v2.base", "disk_size"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "scale_type", "vkcs_kubernetes_node_group_v2.base", "scale_type"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "auto_scale_min_size", "vkcs_kubernetes_node_group_v2.base", "auto_scale_min_size"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "auto_scale_max_size", "vkcs_kubernetes_node_group_v2.base", "auto_scale_max_size"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "auto_scale_node_count", "vkcs_kubernetes_node_group_v2.base", "auto_scale_node_count"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "parallel_upgrade_chunk", "vkcs_kubernetes_node_group_v2.base", "parallel_upgrade_chunk"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "labels.%", "vkcs_kubernetes_node_group_v2.base", "labels.%"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "labels.environment", "vkcs_kubernetes_node_group_v2.base", "labels.environment"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "labels.foo", "vkcs_kubernetes_node_group_v2.base", "labels.foo"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_node_group_v2.node_group", "taints.#", "vkcs_kubernetes_node_group_v2.base", "taints.#"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_node_group_v2.node_group", "uuid"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_node_group_v2.node_group", "created_at"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_node_group_v2.node_group", "region"),
					resource.TestCheckNoResourceAttr("data.vkcs_kubernetes_node_group_v2.node_group", "fixed_scale_node_count"),
				),
			},
		},
	})
}

func testAccCheckKubernetesNodeGroupV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ct, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find node group data source: %s", n)
		}

		if ct.Primary.ID == "" {
			return fmt.Errorf("node group data source ID is not set")
		}

		return nil
	}
}

const testAccKubernetesNodeGroupV2DataSourceClusterBase = `
{{ .TestAccKubernetesNetworkingBase }}

data "vkcs_compute_flavor" "base" {
  name = "Standard-6-12"
}

resource "vkcs_kubernetes_cluster_v2" "base" {
  name                   = "{{ .ClusterName }}"
  version                = "v1.34.2"
  cluster_type           = "standard"
  availability_zones     = ["MS1"]
  master_count           = 1
  master_flavor          = data.vkcs_compute_flavor.base.id
  network_plugin         = "calico"
  pods_ipv4_cidr         = "10.100.0.0/16"
  network_id             = vkcs_networking_network.base.id
  subnet_id              = vkcs_networking_subnet.base.id
  loadbalancer_subnet_id = vkcs_networking_subnet.base.id

  depends_on = [
    vkcs_networking_router_interface.base,
  ]
}
`

const testAccKubernetesNodeGroupV2DataSource = `
{{ .TestAccKubernetesNodeGroupV2DataSourceClusterBase }}

data "vkcs_compute_flavor" "node_flavor" {
  name = "Standard-6-12"
}

resource "vkcs_kubernetes_node_group_v2" "base" {
  cluster_id               = vkcs_kubernetes_cluster_v2.base.id
  name                     = "tfacc-ds-full-v2"
  node_flavor              = data.vkcs_compute_flavor.node_flavor.id
  availability_zone        = "MS1"
  disk_type                = "ceph-ssd"
  disk_size                = 100
  scale_type               = "auto_scale"
  auto_scale_min_size      = 1
  auto_scale_max_size      = 3
  parallel_upgrade_chunk   = 30

  labels = {
    environment = "test"
    foo         = "bar"
  }

  taints = [
    {
      key    = "key1"
      value  = "value1"
      effect = "NoSchedule"
    },
    {
      key    = "key2"
      value  = "value2"
      effect = "PreferNoSchedule"
    }
  ]
}
`

const testAccKubernetesNodeGroupV2DataSourceFullRead = `
{{ .TestAccKubernetesNodeGroupV2DataSourceFull }}

data "vkcs_kubernetes_node_group_v2" "node_group" {
  id = vkcs_kubernetes_node_group_v2.base.id
}
`
