package kubernetes_test

import (
	"errors"
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func TestAccKubernetesNodeGroupV2_basic(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-ng-basic-v2-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesClusterV2Base, map[string]string{
		"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
		"ClusterName":                     clusterName,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupV2Base, map[string]string{
					"TestAccKubernetesNodeGroupV2ClusterBase": clusterConfig,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesNodeGroupV2Exists("vkcs_kubernetes_node_group_v2.base"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group_v2.base", "cluster_id", "vkcs_kubernetes_cluster_v2.base", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.base", "name", "tfacc-basic-v2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.base", "scale_type", "fixed_scale"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.base", "fixed_scale_node_count", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.base", "availability_zone", "MS1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.base", "disk_type", "ceph-ssd"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.base", "disk_size", "30"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.base", "parallel_upgrade_chunk", "20"),

					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group_v2.base", "node_flavor", "data.vkcs_compute_flavor.node_flavor", "id"),

					resource.TestCheckResourceAttrSet("vkcs_kubernetes_node_group_v2.base", "id"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_node_group_v2.base", "uuid"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_node_group_v2.base", "created_at"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_node_group_v2.base", "region"),

					resource.TestCheckNoResourceAttr("vkcs_kubernetes_node_group_v2.base", "auto_scale_min_size"),
					resource.TestCheckNoResourceAttr("vkcs_kubernetes_node_group_v2.base", "auto_scale_max_size"),
					resource.TestCheckNoResourceAttr("vkcs_kubernetes_node_group_v2.base", "auto_scale_node_count"),
					resource.TestCheckNoResourceAttr("vkcs_kubernetes_node_group_v2.base", "labels"),
					resource.TestCheckNoResourceAttr("vkcs_kubernetes_node_group_v2.base", "taints"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_node_group_v2.base"),
		},
	})
}

func TestAccKubernetesNodeGroupV2_full(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-ng-full-v2-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesClusterV2Base, map[string]string{
		"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
		"ClusterName":                     clusterName,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupV2Full, map[string]string{
					"TestAccKubernetesNodeGroupV2ClusterBase": clusterConfig,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesNodeGroupV2Exists("vkcs_kubernetes_node_group_v2.full"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group_v2.full", "cluster_id", "vkcs_kubernetes_cluster_v2.base", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "name", "tfacc-full-v2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "scale_type", "auto_scale"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "auto_scale_min_size", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "auto_scale_max_size", "3"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "auto_scale_node_count", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "availability_zone", "MS1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "disk_type", "ceph-ssd"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "disk_size", "100"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "parallel_upgrade_chunk", "30"),

					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group_v2.full", "node_flavor", "data.vkcs_compute_flavor.node_flavor", "id"),

					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "labels.%", "2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "labels.environment", "test"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "labels.foo", "bar"),

					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "taints.#", "2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "taints.0.key", "key1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "taints.0.value", "value1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "taints.1.key", "key2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "taints.1.value", "value2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.full", "taints.1.effect", "PreferNoSchedule"),

					resource.TestCheckResourceAttrSet("vkcs_kubernetes_node_group_v2.full", "id"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_node_group_v2.full", "uuid"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_node_group_v2.full", "created_at"),
					resource.TestCheckResourceAttrSet("vkcs_kubernetes_node_group_v2.full", "region"),

					resource.TestCheckNoResourceAttr("vkcs_kubernetes_node_group_v2.full", "fixed_scale_node_count"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_node_group_v2.full"),
		},
	})
}

func TestAccKubernetesNodeGroupV2_update(t *testing.T) {
	t.Parallel()

	clusterName := "tfacc-ng-update-v2-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesClusterV2Base, map[string]string{
		"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
		"ClusterName":                     clusterName,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupV2UpdateOld, map[string]string{
					"TestAccKubernetesNodeGroupV2ClusterBase": clusterConfig,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesNodeGroupV2Exists("vkcs_kubernetes_node_group_v2.update"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group_v2.update", "cluster_id", "vkcs_kubernetes_cluster_v2.base", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "name", "tfacc-update-v2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "scale_type", "fixed_scale"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "fixed_scale_node_count", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "availability_zone", "MS1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "disk_type", "ceph-ssd"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "disk_size", "50"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "parallel_upgrade_chunk", "20"),

					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group_v2.update", "node_flavor", "data.vkcs_compute_flavor.node_flavor", "id"),

					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "labels.%", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "labels.environment", "test"),

					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "taints.#", "1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "taints.0.key", "key1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "taints.0.value", "value1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "taints.0.effect", "NoSchedule"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesNodeGroupV2UpdateNew, map[string]string{
					"TestAccKubernetesNetworkingBase":         testAccKubernetesNetworkingBase,
					"TestAccKubernetesNodeGroupV2ClusterBase": clusterConfig,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesNodeGroupV2Exists("vkcs_kubernetes_node_group_v2.update"),
					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group_v2.update", "cluster_id", "vkcs_kubernetes_cluster_v2.base", "id"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "name", "tfacc-update-v2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "scale_type", "auto_scale"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "auto_scale_min_size", "2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "auto_scale_max_size", "5"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "auto_scale_node_count", "2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "availability_zone", "MS1"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "disk_type", "ceph-ssd"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "disk_size", "50"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "parallel_upgrade_chunk", "25"),

					resource.TestCheckResourceAttrPair("vkcs_kubernetes_node_group_v2.update", "node_flavor", "data.vkcs_compute_flavor.node_flavor", "id"),

					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "labels.%", "2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "labels.environment", "production"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "labels.foo", "baz"),

					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "taints.#", "2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "taints.0.key", "key2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "taints.0.value", "value2"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "taints.0.effect", "NoExecute"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "taints.1.key", "key3"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "taints.1.value", "value3"),
					resource.TestCheckResourceAttr("vkcs_kubernetes_node_group_v2.update", "taints.1.effect", "PreferNoSchedule"),
				),
			},
			acctest.ImportStep("vkcs_kubernetes_node_group_v2.update"),
		},
	})
}

func testAccCheckKubernetesNodeGroupV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Kubernetes next-generation node group not found: %s", n)
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

		_, err = nodegroups.GetByID(client, rs.Primary.ID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return errors.New("Kubernetes node group not found")
			}
			return err
		}

		return nil
	}
}

const testAccKubernetesNodeGroupV2Base = `
{{ .TestAccKubernetesNodeGroupV2ClusterBase }}

data "vkcs_compute_flavor" "worker" {
  name = "Standard-6-12"
}

resource "vkcs_kubernetes_node_group_v2" "base" {
  cluster_id               = vkcs_kubernetes_cluster_v2.base.id
  name                     = "tfacc-basic-v2"
  node_flavor              = data.vkcs_compute_flavor.worker.id
  availability_zone        = "MS1"
  disk_type                = "ceph-ssd"
  disk_size                = 30
  scale_type               = "fixed_scale"
  fixed_scale_node_count   = 1
  parallel_upgrade_chunk   = 20
}
`

const testAccKubernetesNodeGroupV2Full = `
{{ .TestAccKubernetesNodeGroupV2ClusterBase }}

data "vkcs_compute_flavor" "node_flavor" {
  name = "Standard-6-12"
}

resource "vkcs_kubernetes_node_group_v2" "full" {
  cluster_id               = vkcs_kubernetes_cluster_v2.base.id
  name                     = "tfacc-full-v2"
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

const testAccKubernetesNodeGroupV2UpdateOld = `
{{ .TestAccKubernetesNodeGroupV2ClusterBase }}

data "vkcs_compute_flavor" "node_flavor" {
  name = "Standard-6-12"
}

resource "vkcs_kubernetes_node_group_v2" "update" {
  cluster_id               = vkcs_kubernetes_cluster_v2.base.id
  name                     = "tfacc-update-v2"
  node_flavor              = data.vkcs_compute_flavor.node_flavor.id
  availability_zone        = "MS1"
  disk_type                = "ceph-ssd"
  disk_size                = 50
  scale_type               = "fixed_scale"
  fixed_scale_node_count   = 1
  parallel_upgrade_chunk   = 20

  labels = {
    environment = "test"
  }

  taints = [
    {
      key    = "key1"
      value  = "value1"
      effect = "NoSchedule"
    }
  ]
}
`

const testAccKubernetesNodeGroupV2UpdateNew = `
{{ .TestAccKubernetesNodeGroupV2ClusterBase }}

data "vkcs_compute_flavor" "node_flavor" {
  name = "Standard-6-18"
}

resource "vkcs_kubernetes_node_group_v2" "update" {
  cluster_id               = vkcs_kubernetes_cluster_v2.base.id
  name                     = "tfacc-update-v2"
  node_flavor              = data.vkcs_compute_flavor.node_flavor.id
  availability_zone        = "MS1"
  disk_type                = "ceph-ssd"
  disk_size                = 50
  scale_type               = "auto_scale"
  auto_scale_min_size      = 2
  auto_scale_max_size      = 5
  parallel_upgrade_chunk   = 25

  labels = {
    environment = "production"
    foo         = "baz"
  }

  taints = [
    {
      key    = "key2"
      value  = "value2"
      effect = "NoExecute"
    },
    {
      key    = "key3"
      value  = "value3"
      effect = "PreferNoSchedule"
    }
  ]
}
`
