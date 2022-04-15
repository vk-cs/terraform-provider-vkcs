package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccKubernetesClusterDataSourceBasic(t *testing.T) {

	var clusterName = "testcluster" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	datasourceName := "data.vkcs_kubernetes_cluster." + clusterName

	createClusterFixture := clusterFixture(clusterName, clusterTemplateID, osFlavorID,
		osKeypairName, osNetworkID, osSubnetworkID, "MS1", 1)

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheckKubernetes(t) },
		ProviderFactories:         testAccProviders,
		CheckDestroy:              testAccCheckKubernetesClusterDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterBasic(createClusterFixture),
			},
			{
				Config: testAccKubernetesClusterDataSourceBasic(
					testAccKubernetesClusterBasic(createClusterFixture), clusterName,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterDataSourceID(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "name", clusterName),
					resource.TestCheckResourceAttr(datasourceName, "master_count", "1"),
					resource.TestCheckResourceAttr(datasourceName, "cluster_template_id", clusterTemplateID),
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

func testAccKubernetesClusterDataSourceBasic(clusterResource, clusterName string) string {
	return fmt.Sprintf(`
%s

data "vkcs_kubernetes_cluster" "`+clusterName+`" {
  name = "${vkcs_kubernetes_cluster.`+clusterName+`.name}"
}
`, clusterResource)
}
