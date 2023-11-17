package kubernetes_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/clusters"
)

func TestAccKubernetesCluster_basic_big(t *testing.T) {
	var cluster clusters.Cluster
	clusterName := "tfacc-basic-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterBasic, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesClusterBase": testAccKubernetesClusterBase, "ClusterName": clusterName}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterExists("vkcs_kubernetes_cluster.basic", &cluster),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster.basic", "name", clusterName),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster.basic", "cluster_template_id", "data.vkcs_kubernetes_clustertemplate.base", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster.basic", "master_count", "1"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_cluster.basic"),
		},
	})
}

func TestAccKubernetesCluster_update_big(t *testing.T) {
	var cluster clusters.Cluster
	clusterName := "tfacc-update-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterUpdate, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesClusterBase": testAccKubernetesClusterBase, "ClusterName": clusterName,
					"EtcdVolumeSize": "10", "CleanVolumes": "false", "KubeLogLevel": "0"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterExists("vkcs_kubernetes_cluster.update", &cluster),
					testAccCheckKubernetesClusterLabels("vkcs_kubernetes_cluster.update", map[string]string{
						"etcd_volume_size": "10",
						"clean_volumes":    "false",
						"kube_log_level":   "0",
					}),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterUpdate, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesClusterBase": testAccKubernetesClusterBase, "ClusterName": clusterName,
					"EtcdVolumeSize": "20", "CleanVolumes": "true", "KubeLogLevel": "8"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterExists("vkcs_kubernetes_cluster.update", &cluster),
					testAccCheckKubernetesClusterLabels("vkcs_kubernetes_cluster.update", map[string]string{
						"etcd_volume_size": "20",
						"clean_volumes":    "true",
						"kube_log_level":   "8",
					}),
				),
			},
		},
	})
}

func TestAccKubernetesCluster_resize_big(t *testing.T) {
	var cluster clusters.Cluster
	clusterName := "tfacc-resize-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterResize, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesClusterBase": testAccKubernetesClusterBase, "ClusterName": clusterName,
					"FlavorName": "Standard-2-8-50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterExists("vkcs_kubernetes_cluster.basic", &cluster),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster.basic", "master_flavor", "data.vkcs_compute_flavor.base", "id"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_cluster.basic"),
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterResize, map[string]string{"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase, "TestAccKubernetesClusterBase": testAccKubernetesClusterBase, "ClusterName": clusterName,
					"FlavorName": "Standard-4-12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterExists("vkcs_kubernetes_cluster.basic", &cluster),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster.basic", "master_flavor", "data.vkcs_compute_flavor.base", "id"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_cluster.basic"),
		},
	})
}

func testAccCheckKubernetesClusterExists(n string, cluster *clusters.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("kubernetes cluster not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("id is not set")
		}

		config, err := clients.ConfigureFromEnv(context.Background())
		if err != nil {
			return fmt.Errorf("Error authenticating clients from environment: %s", err)
		}

		client, err := config.ContainerInfraV1Client(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating Kubernetes API client: %s", err)
		}

		found, err := clusters.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found == nil {
			return errors.New("kubernetes cluster not found")
		}

		*cluster = *found
		return nil
	}
}

func testAccCheckKubernetesClusterLabels(name string, labels map[string]string) resource.TestCheckFunc {
	testFuncs := make([]resource.TestCheckFunc, 0)
	for k, v := range labels {
		testFuncs = append(testFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("labels.%s", k), v))
		testFuncs = append(testFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("all_labels.%s", k), v))
	}
	return resource.ComposeTestCheckFunc(testFuncs...)
}

const testAccKubernetesNetworkingBase = `
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

const testAccKubernetesClusterBase = `
data "vkcs_kubernetes_clustertemplate" "base" {
  version = "1.24"
}
`

const testAccKubernetesClusterBasic = `
{{ .TestAccKubernetesNetworkingBase }}
{{ .TestAccKubernetesClusterBase }}

resource "vkcs_kubernetes_cluster" "basic" {
  name                = "{{ .ClusterName }}"
  cluster_template_id = data.vkcs_kubernetes_clustertemplate.base.id
  master_count        = 1
  network_id          = vkcs_networking_network.base.id
  subnet_id           = vkcs_networking_subnet.base.id
  floating_ip_enabled = false
  availability_zone   = "MS1"

  depends_on = [
    vkcs_networking_router_interface.base,
  ]
}
`

const testAccKubernetesClusterResize = `
{{ .TestAccKubernetesNetworkingBase }}
{{ .TestAccKubernetesClusterBase }}

data "vkcs_compute_flavor" "base" {
  name = "{{ .FlavorName }}"
}

resource "vkcs_kubernetes_cluster" "basic" {
  name                = "{{ .ClusterName }}"
  cluster_template_id = data.vkcs_kubernetes_clustertemplate.base.id
  master_flavor       = data.vkcs_compute_flavor.base.id
  master_count        = 1
  network_id          = vkcs_networking_network.base.id
  subnet_id           = vkcs_networking_subnet.base.id
  floating_ip_enabled = false
  availability_zone   = "MS1"

  depends_on = [
    vkcs_networking_router_interface.base,
  ]
}
`

const testAccKubernetesClusterUpdate = `
{{ .TestAccKubernetesNetworkingBase }}
{{ .TestAccKubernetesClusterBase }}

resource "vkcs_kubernetes_cluster" "update" {
  name                = "{{ .ClusterName }}"
  cluster_template_id = data.vkcs_kubernetes_clustertemplate.base.id
  master_count        = 1
  network_id          = vkcs_networking_network.base.id
  subnet_id           = vkcs_networking_subnet.base.id
  floating_ip_enabled = false
  availability_zone   = "MS1"

  labels = {
	etcd_volume_size = "{{ .EtcdVolumeSize }}"
	clean_volumes    = "{{ .CleanVolumes }}"
	kube_log_level   = "{{ .KubeLogLevel }}"
  }

  depends_on = [
    vkcs_networking_router_interface.base,
  ]
}
`
