package kubernetes_test

import (
	"fmt"
	"testing"

	acctest_helper "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v1/clusters"
)

func TestAccKubernetesClusterDataSource_basic_big(t *testing.T) {
	var cluster clusters.Cluster
	clusterName := "tfacc-basic-" + acctest_helper.RandStringFromCharSet(5, acctest_helper.CharSetAlphaNum)
	clusterConfig := acctest.AccTestRenderConfig(testAccKubernetesClusterBasic, map[string]string{
		"TestAccKubernetesNetworkingBase": testAccKubernetesNetworkingBase,
		"TestAccKubernetesClusterBase":    testAccKubernetesClusterBase,
		"ClusterName":                     clusterName,
	})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: clusterConfig,
				Check:  testAccCheckKubernetesClusterExists("vkcs_kubernetes_cluster.basic", &cluster),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccKubernetesClusterDataSourceBasic, map[string]string{"TestAccKubernetesClusterBasic": clusterConfig, "ClusterName": clusterName}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterDataSourceID("data.vkcs_kubernetes_cluster.cluster"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster.cluster", "name", "vkcs_kubernetes_cluster.basic", "name"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster.cluster", "master_count", "vkcs_kubernetes_cluster.basic", "master_count"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster.cluster", "cluster_template_id", "vkcs_kubernetes_cluster.basic", "cluster_template_id"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster.cluster", "cluster_type", "vkcs_kubernetes_cluster.basic", "cluster_type"),
					resource.TestCheckResourceAttrPair("data.vkcs_kubernetes_cluster.cluster", "availability_zone", "vkcs_kubernetes_cluster.basic", "availability_zone"),
				),
			},
		},
	})
}

func testAccCheckKubernetesClusterDataSourceID(n string) resource.TestCheckFunc {
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

const testAccKubernetesClusterDataSourceBasic = `
{{ .TestAccKubernetesClusterBasic }}

data "vkcs_kubernetes_cluster" "cluster" {
  name = "{{ .ClusterName }}"
}
`
