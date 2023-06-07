package kubernetes_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesAddon_basic_big(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesAddonBasic, map[string]string{
					"TestAccKubernetesAddonNetworkingBase": testAccKubernetesAddonNetworkingBase,
					"TestAccKubernetesAddonClusterBase":    testAccKubernetesAddonClusterBase,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_kubernetes_addon.basic", "namespace", "default"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_addon.basic", "name", "ingress-nginx"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_addon.basic", "cluster_id", "vkcs_kubernetes_cluster.cluster", "id"),
				),
			},
			{
				ResourceName:      "vkcs_kubernetes_addon.basic",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					for _, rs := range s.RootModule().Resources {
						if rs.Type != "vkcs_kubernetes_addon" {
							continue
						}

						id := rs.Primary.Attributes["id"]
						clusterID := rs.Primary.Attributes["cluster_id"]

						return fmt.Sprintf("%s/%s", clusterID, id), nil
					}
					return "", fmt.Errorf("Addon not found")
				},
			},
		},
	})
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
data "vkcs_compute_flavor" "base" {
  name = "Standard-4-12"
}

data "vkcs_kubernetes_clustertemplate" "ct" {
  version = "1.24"
}

resource "vkcs_kubernetes_cluster" "cluster" {
  depends_on = [
    vkcs_networking_router_interface.base,
  ]

  name                = "tfacc-cluster"
  cluster_template_id = data.vkcs_kubernetes_clustertemplate.ct.id
  master_flavor       = data.vkcs_compute_flavor.base.id
  master_count        = 1

  network_id          = vkcs_networking_network.base.id
  subnet_id           = vkcs_networking_subnet.base.id
  floating_ip_enabled = true
  availability_zone   = "MS1"
  insecure_registries = ["1.2.3.4"]
}

resource "vkcs_kubernetes_node_group" "default-ng" {
  cluster_id = vkcs_kubernetes_cluster.cluster.id
  name       = "tfacc-default-ng"
  node_count = 1
  max_nodes  = 5
  min_nodes  = 1
}
`

const testAccKubernetesAddonBasic = `
{{ .TestAccKubernetesAddonNetworkingBase }}
{{ .TestAccKubernetesAddonClusterBase }}

data "vkcs_kubernetes_addon" "ingress-nginx" {
  cluster_id = vkcs_kubernetes_cluster.cluster.id
  name       = "ingress-nginx"
  version    = "4.1.4"
}

resource "vkcs_kubernetes_addon" "basic" {
  cluster_id           = vkcs_kubernetes_cluster.cluster.id
  addon_id             = data.vkcs_kubernetes_addon.ingress-nginx.id
  configuration_values = data.vkcs_kubernetes_addon.ingress-nginx.configuration_values
  depends_on           = [vkcs_kubernetes_node_group.default-ng]
}
`
