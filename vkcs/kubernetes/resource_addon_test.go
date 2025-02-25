package kubernetes_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/clusters"
)

func TestAccKubernetesAddon_basic_big(t *testing.T) {
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
				Config: acctest.AccTestRenderConfig(testAccKubernetesAddonBasic, map[string]string{"TestAccKubernetesAddonClusterBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_addon.addon", "addon_id", "data.vkcs_kubernetes_addon.ingress-nginx", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_addon.addon", "cluster_id", "vkcs_kubernetes_cluster.cluster", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_addon.addon", "namespace", "ingress-nginx"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_addon.addon", "name", "ingress-nginx"),
				),
			},
			{
				ResourceName:            "vkcs_kubernetes_addon.addon",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccKubernetesAddonImportStateID,
				ImportStateVerifyIgnore: []string{"configuration_values"},
			},
		},
	})
}

func TestAccKubernetesAddon_full_big(t *testing.T) {
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
				Config: acctest.AccTestRenderConfig(testAccKubernetesAddonFull, map[string]string{"TestAccKubernetesAddonClusterBase": baseConfig}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_addon.addon", "addon_id", "data.vkcs_kubernetes_addon.ingress-nginx", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_addon.addon", "cluster_id", "vkcs_kubernetes_cluster.cluster", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_addon.addon", "name", "ingress-nginx"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_addon.addon", "namespace", "ingress-nginx"),
				),
			},
			{
				ResourceName:            "vkcs_kubernetes_addon.addon",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccKubernetesAddonImportStateID,
				ImportStateVerifyIgnore: []string{"configuration_values"},
			},
		},
	})
}

func testAccKubernetesAddonImportStateID(s *terraform.State) (string, error) {
	for name, rs := range s.RootModule().Resources {
		if name != "vkcs_kubernetes_addon.addon" {
			continue
		}

		id := rs.Primary.Attributes["id"]
		clusterID := rs.Primary.Attributes["cluster_id"]

		return fmt.Sprintf("%s/%s", clusterID, id), nil
	}
	return "", fmt.Errorf("addon not found")
}

func uniqueKubernetesNodeGroupName() string {
	t := time.Now()
	return fmt.Sprintf("tfacc-def-ng-%dh-%dm", t.Hour(), t.Minute())
}

const testAccKubernetesAddonNetworkingBase = `
resource "vkcs_networking_network" "base" {
  name           = "tfacc-base-net"
  admin_state_up = true
}

resource "vkcs_networking_subnet" "base" {
  name            = "tfacc-base-subnet"
  network_id      = vkcs_networking_network.base.id
  cidr            = "192.168.199.0/24"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

data "vkcs_networking_network" "base-extnet" {
  name = "ext-net"
}

resource "vkcs_networking_router" "base" {
  name                = "tfacc-base-router"
  admin_state_up      = true
  external_network_id = data.vkcs_networking_network.base-extnet.id
}

resource "vkcs_networking_router_interface" "base" {
  router_id = vkcs_networking_router.base.id
  subnet_id = vkcs_networking_subnet.base.id
}
`

const testAccKubernetesAddonClusterBase = `
{{ .TestAccKubernetesAddonNetworkingBase }}

data "vkcs_compute_flavor" "base" {
  name = "Standard-4-12"
}

data "vkcs_kubernetes_clustertemplate" "ct" {
  version = "1.31"
}

resource "vkcs_kubernetes_cluster" "cluster" {
  name                = "tfacc-cls-{{.Suffix}}"
  cluster_template_id = data.vkcs_kubernetes_clustertemplate.ct.id
  master_flavor       = data.vkcs_compute_flavor.base.id
  master_count        = 1

  network_id          = vkcs_networking_network.base.id
  subnet_id           = vkcs_networking_subnet.base.id
  floating_ip_enabled = true
  availability_zone   = "MS1"
  insecure_registries = ["1.2.3.4"]

  depends_on = [
    vkcs_networking_router_interface.base,
  ]
}

resource "vkcs_kubernetes_node_group" "default-ng" {
  cluster_id = vkcs_kubernetes_cluster.cluster.id
  name       = "{{.NodeGroupName}}"
  node_count = 1
  max_nodes  = 5
  min_nodes  = 1
}
`

const testAccKubernetesAddonBasic = `
{{ .TestAccKubernetesAddonClusterBase }}

data "vkcs_kubernetes_addon" "ingress-nginx" {
  cluster_id = vkcs_kubernetes_cluster.cluster.id
  name       = "ingress-nginx"
  version    = "4.12.0"
}

resource "vkcs_kubernetes_addon" "addon" {
  cluster_id           = vkcs_kubernetes_cluster.cluster.id
  addon_id             = data.vkcs_kubernetes_addon.ingress-nginx.id
  namespace            = "ingress-nginx"
}
`

const testAccKubernetesAddonFull = `
{{ .TestAccKubernetesAddonClusterBase }}

data "vkcs_kubernetes_addon" "ingress-nginx" {
  cluster_id = vkcs_kubernetes_cluster.cluster.id
  name       = "ingress-nginx"
  version    = "4.12.0"
}

resource "vkcs_kubernetes_addon" "addon" {
  cluster_id           = vkcs_kubernetes_cluster.cluster.id
  addon_id             = data.vkcs_kubernetes_addon.ingress-nginx.id
  name                 = "ingress-nginx"
  namespace            = "ingress-nginx"
  configuration_values = data.vkcs_kubernetes_addon.ingress-nginx.configuration_values
}
`
