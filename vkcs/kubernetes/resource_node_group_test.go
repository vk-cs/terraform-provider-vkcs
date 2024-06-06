package kubernetes_test

import (
	"errors"
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/nodegroups"
)

func TestAccKubernetesNodeGroup_basic_big(t *testing.T) {
	var nodeGroup nodegroups.NodeGroup
	clusterName := "tfacc-ng-basic-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesNodeGroupClusterBase, map[string]string{"ClusterName": clusterName})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupBasic, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesNodeGroupClusterBase": clusterConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesNodeGroupExists("vkcs_kubernetes_node_group.basic", &nodeGroup),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group.basic", "cluster_id", "vkcs_kubernetes_cluster.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group.basic", "flavor_id", "data.vkcs_compute_flavor.base", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.basic", "node_count", "1"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_node_group.basic"),
		},
	})
}

func TestAccKubernetesNodeGroup_resize_big(t *testing.T) {
	var nodeGroup nodegroups.NodeGroup
	clusterName := "tfacc-ng-resize-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesNodeGroupClusterBase, map[string]string{"ClusterName": clusterName})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupResize, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesNodeGroupClusterBase": clusterConfig,
					"FlavorName": "Standard-2-8-50"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesNodeGroupExists("vkcs_kubernetes_node_group.resize", &nodeGroup),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group.resize", "flavor_id", "data.vkcs_compute_flavor.resize", "id"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupResize, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesNodeGroupClusterBase": clusterConfig,
					"FlavorName": "Standard-4-12"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesNodeGroupExists("vkcs_kubernetes_node_group.resize", &nodeGroup),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group.resize", "flavor_id", "data.vkcs_compute_flavor.resize", "id"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_node_group.resize"),
		},
	})
}

func testAccCheckKubernetesNodeGroupExists(n string, nodeGroup *nodegroups.NodeGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("kubernetes node group not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("id is not set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		client, err := config.ContainerInfraV1Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating Kubernetes API client: %s", err)
		}

		found, err := nodegroups.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found == nil {
			return errors.New("kubernetes node group not found")
		}

		*nodeGroup = *found
		return nil
	}
}

const testAccKubernetesNodeGroupClusterBase = `
data "vkcs_compute_flavor" "base" {
  name = "Standard-2-8-50"
}

data "vkcs_kubernetes_clustertemplate" "base" {
  version = "1.27"
}

resource "vkcs_kubernetes_cluster" "base" {
  name                = "{{ .ClusterName }}"
  cluster_template_id = data.vkcs_kubernetes_clustertemplate.base.id
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
`

const testAccKubernetesNodeGroupBasic = `
{{ .TestAccKubernetesNetworkingBase }}
{{ .TestAccKubernetesNodeGroupClusterBase }}

resource "vkcs_kubernetes_node_group" "basic" {
  cluster_id          = vkcs_kubernetes_cluster.base.id
  name                = "tfacc-basic"
  flavor_id           = data.vkcs_compute_flavor.base.id
  node_count          = 1
  max_nodes           = 5
  min_nodes           = 1
  autoscaling_enabled = false
}
`

const testAccKubernetesNodeGroupResize = `
{{ .TestAccKubernetesNetworkingBase }}
{{ .TestAccKubernetesNodeGroupClusterBase }}

data "vkcs_compute_flavor" "resize" {
  name = "{{ .FlavorName }}"
}

resource "vkcs_kubernetes_node_group" "resize" {
  cluster_id          = vkcs_kubernetes_cluster.base.id
  name                = "tfacc-resize"
  flavor_id           = data.vkcs_compute_flavor.resize.id
  node_count          = 1
  max_nodes           = 5
  min_nodes           = 1
  autoscaling_enabled = false
}
`
