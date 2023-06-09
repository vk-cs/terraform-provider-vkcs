package kubernetes_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesClusterTemplatesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterTemplatesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccKubernetesClusterTemplatesCheckNotEmpty("data.vkcs_kubernetes_clustertemplates.templates"),
				),
			},
		},
	})
}

func testAccKubernetesClusterTemplatesCheckNotEmpty(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("root module has no resource called %s", resourceName)
		}

		templates, ok := rs.Primary.Attributes["cluster_templates.#"]
		if !ok {
			return fmt.Errorf("cluster_templates attribute is missing.")
		}

		templatesQuantity, err := strconv.Atoi(templates)
		if err != nil {
			return fmt.Errorf("error parsing templates (%s) into integer: %s", templates, err)
		}

		if templatesQuantity == 0 {
			return fmt.Errorf("No templates found, this is probably a bug.")
		}

		return nil
	}
}

const testAccKubernetesClusterTemplatesDataSourceBasic = `
data "vkcs_kubernetes_clustertemplates" "templates" {}
`
