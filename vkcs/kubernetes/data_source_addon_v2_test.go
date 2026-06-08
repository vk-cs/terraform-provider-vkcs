package kubernetes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccKubernetesAddonV2DataSource_byNameAndVersion(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAddonV2DataSourceByNameConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_addon_v2.by_name", "name", "ingress-nginx"),
					resource.TestCheckResourceAttr("data.vkcs_kubernetes_addon_v2.by_name", "version", "4.12.1"),

					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_addon_v2.by_name", "id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_addon_v2.by_name", "addon_id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_addon_v2.by_name", "region"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_addon_v2.by_name", "values_template"),
				),
			},
		},
	})
}

func TestAccKubernetesAddonV2DataSource_byID(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAddonV2DataSourceByIDConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.vkcs_kubernetes_addon_v2.by_id", "id",
						"data.vkcs_kubernetes_addon_v2.by_name", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_kubernetes_addon_v2.by_id", "addon_id",
						"data.vkcs_kubernetes_addon_v2.by_name", "addon_id",
					),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_kubernetes_addon_v2.by_id", "name",
						"data.vkcs_kubernetes_addon_v2.by_name", "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_kubernetes_addon_v2.by_id", "version",
						"data.vkcs_kubernetes_addon_v2.by_name", "version",
					),

					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_addon_v2.by_id", "id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_addon_v2.by_id", "addon_id"),
					resource.TestCheckResourceAttrSet("data.vkcs_kubernetes_addon_v2.by_id", "region"),
				),
			},
		},
	})
}

const testAccKubernetesAddonV2DataSourceByNameConfig = `
data "vkcs_kubernetes_addon_v2" "by_name" {
  name    = "ingress-nginx"
  version = "4.12.1"
}
`

const testAccKubernetesAddonV2DataSourceByIDConfig = `
data "vkcs_kubernetes_addon_v2" "by_name" {
  name    = "ingress-nginx"
  version = "4.12.1"
}

data "vkcs_kubernetes_addon_v2" "by_id" {
  id       = data.vkcs_kubernetes_addon_v2.by_name.id
}
`
