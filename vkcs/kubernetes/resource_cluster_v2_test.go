package kubernetes_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
)

func TestAccKubernetesClusterV2_basic(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-cl-basic-v2-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterV2Basic, map[string]string{
					"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
					"ClusterName":                     clusterName,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterV2Exists("vkcs_kubernetes_cluster_v2.basic"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "name", clusterName),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "version", "v1.34.2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "cluster_type", "standard"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "availability_zones.#", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "availability_zones.0", "MS1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "master_count", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "network_plugin", "calico"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "pods_ipv4_cidr", "10.100.0.0/16"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "public_ip", "false"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "master_disks.#", "3"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "description", ""),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "api_lb_fip", ""),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "insecure_registries.#", "0"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "labels.#", "0"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "loadbalancer_allowed_cidrs.#", "0"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.basic", "node_groups.#", "0"),

					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.basic", "master_flavor", "data.vkcs_compute_flavor.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.basic", "network_id", "vkcs_networking_network.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.basic", "subnet_id", "vkcs_networking_subnet.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.basic", "loadbalancer_subnet_id", "vkcs_networking_subnet.base", "id"),

					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.basic", "id"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.basic", "uuid"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.basic", "external_network_id"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.basic", "k8s_config"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.basic", "created_at"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.basic", "status"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.basic", "project_id"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.basic", "api_address"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.basic", "api_lb_vip"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.basic", "region"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_cluster_v2.basic"),
		},
	})
}

func TestAccKubernetesClusterV2_full(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-cl-full-v2-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterUUID := uuid.NewString()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterV2Full, map[string]string{
					"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
					"ClusterName":                     clusterName,
					"ClusterUUID":                     clusterUUID,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterV2Exists("vkcs_kubernetes_cluster_v2.full"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "name", clusterName),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "uuid", clusterUUID),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "description", "Test cluster with all attributes"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "version", "v1.34.2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "cluster_type", "standard"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "availability_zones.#", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "availability_zones.0", "MS1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "master_count", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "network_plugin", "calico"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "pods_ipv4_cidr", "10.100.0.0/16"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "public_ip", "true"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "loadbalancer_allowed_cidrs.#", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "loadbalancer_allowed_cidrs.0", "10.0.0.0/8"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "insecure_registries.#", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "insecure_registries.0", "registry.example.com:5000"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "labels.%", "2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "labels.foo", "bar"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "labels.environment", "test"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "master_disks.#", "3"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.full", "node_groups.#", "0"),

					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.full", "external_network_id", "data.vkcs_networking_network.base-extnet", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.full", "master_flavor", "data.vkcs_compute_flavor.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.full", "network_id", "vkcs_networking_network.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.full", "subnet_id", "vkcs_networking_subnet.base", "id"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.full", "loadbalancer_subnet_id", "vkcs_networking_subnet.base", "id"),

					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.full", "id"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.full", "k8s_config"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.full", "created_at"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.full", "status"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.full", "project_id"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.full", "api_lb_fip"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.full", "api_lb_vip"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.full", "api_address"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_cluster_v2.full", "region"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_cluster_v2.full"),
		},
	})
}

func TestAccKubernetesClusterV2_scale(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-cl-scale-v2-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterV2Scale, map[string]string{
					"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
					"ClusterName":                     clusterName,
					"FlavorName":                      "Standard-6-12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterV2Exists("vkcs_kubernetes_cluster_v2.scale"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.scale", "name", clusterName),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.scale", "master_flavor", "data.vkcs_compute_flavor.base", "id"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterV2Scale, map[string]string{
					"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
					"ClusterName":                     clusterName,
					"FlavorName":                      "Standard-6-18",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterV2Exists("vkcs_kubernetes_cluster_v2.scale"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.scale", "name", clusterName),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_cluster_v2.scale", "master_flavor", "data.vkcs_compute_flavor.base", "id"),
				),
			},
		},
	})
}

/*
func TestAccKubernetesClusterV2_upgrade(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-cl-upgrade-v2-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterV2Upgrade, map[string]string{
					"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
					"ClusterName":                     clusterName,
					"ClusterVersion":                  "v1.32.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterV2Exists("vkcs_kubernetes_cluster_v2.upgrade"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.upgrade", "name", clusterName),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.upgrade", "version", "v1.32.1"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterV2Upgrade, map[string]string{
					"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
					"ClusterName":                     clusterName,
					"ClusterVersion":                  "v1.33.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterV2Exists("vkcs_kubernetes_cluster_v2.upgrade"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.upgrade", "name", clusterName),
					resource.TestCheckResourceAttr("vkcs_kubernetes_cluster_v2.upgrade", "version", "v1.33.3"),
				),
			},
		},
	})
}
*/

func testAccCheckKubernetesClusterV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Kubernetes next-generation cluster not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("id is not set")
		}

		opts := clients.ConfigOpts{}
		config, err := opts.LoadAndValidate()
		if err != nil {
			return fmt.Errorf("Error authenticating clients from environment: %s", err)
		}

		client, err := config.ManagedK8SClient(acctest.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating Kubernetes API client: %s", err)
		}

		found, err := clusters.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found == nil {
			return errors.New("Kubernetes cluster not found")
		}

		return nil
	}
}

const testAccKubernetesClusterV2Basic = `
{{ .TestAccKubernetesNetworkingBase }}

data "vkcs_compute_flavor" "base" {
  name = "Standard-6-12"
}

resource "vkcs_kubernetes_cluster_v2" "basic" {
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

const testAccKubernetesClusterV2Full = `
{{ .TestAccKubernetesNetworkingBase }}

data "vkcs_compute_flavor" "base" {
  name = "Standard-6-12"
}

resource "vkcs_kubernetes_cluster_v2" "full" {
  name                   = "{{ .ClusterName }}"
  uuid                   = "{{ .ClusterUUID }}"
  description            = "Test cluster with all attributes"
  version                = "v1.34.2"
  cluster_type           = "standard"
  availability_zones     = ["MS1"]
  master_count           = 1
  master_flavor          = data.vkcs_compute_flavor.base.id
  network_plugin         = "calico"
  pods_ipv4_cidr         = "10.100.0.0/16"
  network_id             = vkcs_networking_network.base.id
  subnet_id              = vkcs_networking_subnet.base.id
  public_ip              = true
  external_network_id    = data.vkcs_networking_network.base-extnet.id
  loadbalancer_subnet_id = vkcs_networking_subnet.base.id
  loadbalancer_allowed_cidrs = ["10.0.0.0/8"]

  insecure_registries = ["registry.example.com:5000"]

  labels = {
    foo         = "bar"
    environment = "test"
  }

  depends_on = [
    vkcs_networking_router_interface.base,
  ]
}
`

const testAccKubernetesClusterV2Scale = `
{{ .TestAccKubernetesNetworkingBase }}

data "vkcs_compute_flavor" "base" {
  name = "{{ .FlavorName }}"
}

resource "vkcs_kubernetes_cluster_v2" "scale" {
  name                   = "{{ .ClusterName }}"
  version                = "v1.34.2"
  cluster_type           = "standard"
  availability_zones     = ["MS1"]
  master_count           = 1
  master_flavor          = data.vkcs_compute_flavor.base.id
  network_plugin         = "calico"
  network_id             = vkcs_networking_network.base.id
  subnet_id              = vkcs_networking_subnet.base.id
  loadbalancer_subnet_id = vkcs_networking_subnet.base.id
  description            = "Cluster for update"
  pods_ipv4_cidr         = "10.100.0.0/16"

  labels = {
    foo = "baz"
    new = "value"
  }

  loadbalancer_allowed_cidrs = ["10.0.0.0/8", "192.168.0.0/16"]

  depends_on = [
    vkcs_networking_router_interface.base,
  ]
}
`

/*
const testAccKubernetesClusterV2Upgrade = `
{{ .TestAccKubernetesNetworkingBase }}

data "vkcs_compute_flavor" "base" {
  name = "Standard-6-12"
}

resource "vkcs_kubernetes_cluster_v2" "upgrade" {
  name                   = "{{ .ClusterName }}"
  version                = "{{ .ClusterVersion }}"
  cluster_type           = "standard"
  availability_zones     = ["MS1"]
  master_count           = 1
  master_flavor          = data.vkcs_compute_flavor.base.id
  network_plugin         = "calico"
  network_id             = vkcs_networking_network.base.id
  subnet_id              = vkcs_networking_subnet.base.id
  loadbalancer_subnet_id = vkcs_networking_subnet.base.id
  description            = "Cluster for update"
  pods_ipv4_cidr         = "10.100.0.0/16"

  labels = {
    foo = "baz"
    new = "value"
  }

  loadbalancer_allowed_cidrs = ["10.0.0.0/8", "192.168.0.0/16"]

  depends_on = [
    vkcs_networking_router_interface.base,
  ]
}
`
*/
