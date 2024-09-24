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

func TestAccKubernetesNodeGroup_fullUpdate_big(t *testing.T) {
	var nodeGroup nodegroups.NodeGroup
	clusterName := "tfacc-ng-basic-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesNodeGroupClusterBase, map[string]string{"ClusterName": clusterName})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupFullUpdateOld, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesNodeGroupClusterBase": clusterConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesNodeGroupExists("vkcs_kubernetes_node_group.full", &nodeGroup),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group.full", "cluster_id", "vkcs_kubernetes_cluster.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group.full", "flavor_id", "data.vkcs_compute_flavor.node_flavor", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "availability_zones.#", "3"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "node_count", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "max_nodes", "5"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "min_nodes", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "autoscaling_enabled", "false"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "labels.#", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "taints.#", "2"),
				),
			},
			{
				ResourceName:            "vkcs_kubernetes_node_group.full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "taints"},
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupFullUpdateNew, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesNodeGroupClusterBase": clusterConfig}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesNodeGroupExists("vkcs_kubernetes_node_group.full", &nodeGroup),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group.full", "cluster_id", "vkcs_kubernetes_cluster.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group.full", "flavor_id", "data.vkcs_compute_flavor.node_flavor", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "availability_zones.#", "3"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "max_nodes", "10"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "min_nodes", "2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "autoscaling_enabled", "true"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "labels.#", "2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.full", "taints.#", "1"),
				),
			},
			{
				ResourceName:            "vkcs_kubernetes_node_group.full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "taints"},
			},
		},
	})
}

func TestAccKubernetesNodeGroup_scale_big(t *testing.T) {
	var nodeGroup nodegroups.NodeGroup
	clusterName := "tfacc-ng-scale-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesNodeGroupClusterBase, map[string]string{"ClusterName": clusterName})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupScale, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesNodeGroupClusterBase": clusterConfig,
					"NodeCount": "1"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesNodeGroupExists("vkcs_kubernetes_node_group.scale", &nodeGroup),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group.scale", "flavor_id", "data.vkcs_compute_flavor.base", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.scale", "node_count", "1"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_node_group.scale"),
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupScale, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesNodeGroupClusterBase": clusterConfig,
					"NodeCount": "2"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesNodeGroupExists("vkcs_kubernetes_node_group.scale", &nodeGroup),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group.scale", "flavor_id", "data.vkcs_compute_flavor.base", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group.scale", "node_count", "2"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_node_group.scale"),
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
  master_flavor       = data.vkcs_compute_flavor.base.id
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

const testAccKubernetesNodeGroupScale = `
{{ .TestAccKubernetesNetworkingBase }}
{{ .TestAccKubernetesNodeGroupClusterBase }}

resource "vkcs_kubernetes_node_group" "scale" {
  cluster_id          = vkcs_kubernetes_cluster.base.id
  name                = "tfacc-scale"
  flavor_id           = data.vkcs_compute_flavor.base.id
  node_count          = {{ .NodeCount }}
  max_nodes           = 5
  min_nodes           = 1
  autoscaling_enabled = false
}
`

const testAccKubernetesNodeGroupFullUpdateOld = `
{{ .TestAccKubernetesNetworkingBase }}
{{ .TestAccKubernetesNodeGroupClusterBase }}

data "vkcs_compute_flavor" "node_flavor" {
  name = "STD2-2-8"
}

resource "vkcs_kubernetes_node_group" "full" {
  cluster_id          = vkcs_kubernetes_cluster.base.id
  name                = "tfacc-full-update"
  flavor_id           = data.vkcs_compute_flavor.node_flavor.id
  availability_zones  = ["ME1", "GZ1", "MS1"]
  node_count          = 1
  max_nodes           = 5
  min_nodes           = 1
  autoscaling_enabled = false
  labels {
    key   = "label1"
    value = "test1"
  }
  taints {
    key    = "taint1"
    value  = "test1"
    effect = "PreferNoSchedule"
  }
  taints {
    key    = "taint2"
    value  = "test2"
    effect = "NoSchedule"
  }
}
`

const testAccKubernetesNodeGroupFullUpdateNew = `
{{ .TestAccKubernetesNetworkingBase }}
{{ .TestAccKubernetesNodeGroupClusterBase }}

data "vkcs_compute_flavor" "node_flavor" {
  name = "STD3-4-12"
}

resource "vkcs_kubernetes_node_group" "full" {
  cluster_id           = vkcs_kubernetes_cluster.base.id
  name                 = "tfacc-full-update"
  flavor_id            = data.vkcs_compute_flavor.node_flavor.id
  availability_zones   = ["ME1", "GZ1", "MS1"]
  node_count           = 2
  max_nodes            = 10
  min_nodes            = 2
  max_node_unavailable = 1
  autoscaling_enabled  = true
  labels {
    key   = "label2"
    value = "test3"
  }
  labels {
    key   = "label3"
    value = "test3"
  }
  taints {
    key    = "taint3"
    value  = "test3"
    effect = "NoExecute"
  }
}
`
