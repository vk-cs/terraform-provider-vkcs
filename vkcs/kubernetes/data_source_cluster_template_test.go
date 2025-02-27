package kubernetes_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesClusterTemplateDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterTemplateDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterTemplateDataSourceID("data.vkcs_kubernetes_clustertemplate.template"),
				),
			},
		},
	})
}

func TestAccKubernetesClusterTemplateDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterTemplateDataSourceQueryUUID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKubernetesClusterTemplateDataSourceID("data.vkcs_kubernetes_clustertemplate.template"),
				),
			},
		},
	})
}

func testAccCheckKubernetesClusterTemplateDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find cluster template data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("cluster template data source ID is not set")
		}

		return nil
	}
}

const testAccKubernetesClusterTemplateDataSourceBasic = `
data "vkcs_kubernetes_clustertemplate" "template" {
	version = "1.31"
}
`

const testAccKubernetesClusterTemplateDataSourceQueryUUID = `
data "vkcs_kubernetes_clustertemplate" "template" {
	id = "42777b6f-e9ea-4fc4-899a-eef08b8b0380"
}
`
